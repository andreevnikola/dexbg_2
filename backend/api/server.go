package api

import (
	"dexbg/common"
	"dexbg/storage"
	"dexbg/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	listenAddr string
	store      storage.Storage
}

func NewServer(listenAddr string, store storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Start() error {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"}
	router.Use(cors.New(config))

	router.POST("/users/join", s.handleAddNewUser)
	router.POST("/users/signin", s.handleLoginToAccount)
	router.GET("/users/authenticate/:key", s.handleAuthenticate)
	return router.Run("localhost:" + s.listenAddr)
}

func (s *Server) handleAddNewUser(c *gin.Context) {
	readyStateError := make(chan int)

	var newUser types.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(500, nil)
		return
	}

	go func() {
		var isUserValid bool = types.ValidateUser(&types.User{
			Username: newUser.Username,
			Password: newUser.Password,
			Phone:    newUser.Phone,
			Mail:     newUser.Mail,
			Fullname: newUser.Fullname,
		})
		if !isUserValid {
			readyStateError <- 406
			return
		}
		readyStateError <- 0
	}()

	go func() {
		user, err := s.store.FindOne(&bson.M{
			"$or": []bson.M{
				{"username": newUser.Username},
				{"mail": newUser.Mail},
			},
		}, "users")
		if err != nil {
			readyStateError <- 500
			return
		}
		if user != nil {
			var errorCode int
			if user["username"] == newUser.Username {
				errorCode = 422
			} else {
				errorCode = 409
			}
			readyStateError <- errorCode
			return
		}
		readyStateError <- 0
	}()

	var hashedPas string
	go func() {
		bytes, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 6)
		if err != nil {
			readyStateError <- 500
			return
		}
		hashedPas = string(bytes)
		readyStateError <- 0
	}()

	var gender int8 = 0
	go func() {
		type Response struct {
			Count       int     `json: "count"`
			Gender      string  `json: "gender"`
			Name        string  `json: "name"`
			Probability float32 `json: "probability"`
		}
		resp, err := http.Get("https://api.genderize.io/?name[]=" + strings.Split(newUser.Fullname, " ")[0] + "&name[]=" + strings.Split(newUser.Fullname, " ")[1])
		if err != nil {
			readyStateError <- 0
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			readyStateError <- 0
			return
		}
		var result []Response
		if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
			readyStateError <- 0
			return
		}
		// fmt.Println("BRUH", result)
		if result[0].Gender == result[1].Gender && result[0].Probability+result[1].Probability > 1.25 {
			if result[0].Gender == "female" {
				gender = 2
			} else {
				gender = 1
			}
		}
		readyStateError <- 0
	}()

	for i := 0; i < 4; i++ {
		var stateError int = <-readyStateError
		if stateError != 0 {
			c.JSON(stateError, nil)
			return
		}
	}

	var key types.Key = common.GenerateKey()

	var data interface{} = bson.D{
		{Key: "username", Value: newUser.Username},
		{Key: "mail", Value: newUser.Mail},
		{Key: "phone", Value: newUser.Phone},
		{Key: "fullname", Value: newUser.Fullname},
		{Key: "password", Value: hashedPas},
		{Key: "gender", Value: gender},
		{Key: "keys", Value: []types.Key{key}},
	}
	res, err := s.store.InsertOne(data, "users")
	if err != nil {
		c.JSON(500, nil)
		return
	}

	c.JSON(200, gin.H{
		"id":     res.InsertedID,
		"key":    key.Key,
		"gender": gender,
	})

}

func (s *Server) handleLoginToAccount(c *gin.Context) {
	var userLogin types.UserLogin
	if err := c.BindJSON(&userLogin); err != nil {
		c.JSON(500, nil)
		return
	}
	user, err := s.store.FindOne(&bson.M{
		"$or": []bson.M{
			{"username": userLogin.Login},
			{"mail": userLogin.Login},
		},
	}, "users")
	if err != nil {
		c.JSON(500, nil)
		return
	}
	if user == nil {
		c.JSON(404, nil)
		return
	}
	accPass, ok := user["password"].(string)
	if !ok {
		fmt.Println(accPass)
		c.JSON(404, nil)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(accPass), []byte(userLogin.Password)); err != nil {
		c.JSON(404, nil)
		return
	}

	var keys reflect.Value = reflect.ValueOf(user["keys"])
	var key types.Key = common.GenerateKey()

	if keys.Len() > 3 {
		if _, err = s.store.Update(&bson.M{
			"$or": []bson.M{
				{"username": userLogin.Login},
				{"mail": userLogin.Login},
			},
		},
			bson.M{
				"$unset": bson.M{"keys.0": 1},
			}, "users"); err != nil {
			c.JSON(500, nil)
			return
		}
		if _, err = s.store.Update(&bson.M{
			"$or": []bson.M{
				{"username": userLogin.Login},
				{"mail": userLogin.Login},
			},
		},
			bson.M{
				"$pull": bson.M{"keys": nil},
			}, "users"); err != nil {
			c.JSON(500, nil)
			return
		}
	}
	if _, err = s.store.Update(&bson.M{
		"$or": []bson.M{
			{"username": userLogin.Login},
			{"mail": userLogin.Login},
		},
	},
		bson.M{
			"$push": bson.M{"keys": key},
		}, "users"); err != nil {
		c.JSON(500, nil)
		return
	}

	c.JSON(200, gin.H{
		"id":       user["_id"].(primitive.ObjectID).Hex(),
		"key":      key.Key,
		"mail":     user["mail"],
		"phone":    user["phone"],
		"username": user["username"],
		"fullname": user["fullname"],
		"gender":   user["gender"],
	})
}

func (s *Server) handleAuthenticate(c *gin.Context) {
	var key string = c.Param("key")
	user, err := s.store.FindOne(&bson.M{
		"keys.key": key,
	}, "users")
	if err != nil {
		c.JSON(500, nil)
		return
	}
	if user == nil {
		c.JSON(401, nil)
		return
	}
	now := time.Now()
	nowInSeconds := int(now.Unix())
	keys, _ := user["keys"].(primitive.A)
	for i, v := range keys {
		lookKey, _ := v.(primitive.M)
		if lookKey["key"] == key && lookKey["expiration"].(int32) > int32(nowInSeconds) {
			var index string = "keys." + strconv.Itoa(i) + ".expiration"
			var updated *mongo.UpdateResult
			if updated, err = s.store.Update(&bson.M{
				"keys.key": key,
			},
				bson.M{
					"$set": bson.M{index: nowInSeconds + 345600},
				}, "users"); err != nil {
				fmt.Println(err)
				c.JSON(500, nil)
				return
			}
			if updated.ModifiedCount == 0 {
				c.JSON(401, nil)
			}
			c.JSON(200, gin.H{
				"id":       user["_id"].(primitive.ObjectID).Hex(),
				"mail":     user["mail"],
				"phone":    user["phone"],
				"username": user["username"],
				"fullname": user["fullname"],
				"gender":   user["gender"],
			})
			return
		}
	}
	c.JSON(401, nil)
}
