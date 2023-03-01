package storage

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage interface {
	InsertOne(data interface{}, collectionString string) (*mongo.InsertOneResult, error)
	InsertMany(data []interface{}, collectionString string) (*mongo.InsertManyResult, error)
	Find(search *bson.D, collectionString string) ([]map[string]interface{}, error)
	FindOne(search *bson.M, collectionString string) (map[string]interface{}, error)
	Update(search *bson.M, update bson.M, collectionString string) (*mongo.UpdateResult, error)
}
