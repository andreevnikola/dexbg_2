package types

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Mail     string `json:"mail"`
	Fullname string `json:"fullname"`
}

type UserLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Key struct {
	Key        string
	Expiration int
}

func ValidateUser(u *User) bool {
	if len([]rune(u.Username)) < 5 || len([]rune(u.Username)) > 20 {
		return false
	}
	if len([]rune(u.Password)) < 5 || len([]rune(u.Password)) > 50 {
		return false
	}
	if len([]rune(u.Phone)) != 12 && len([]rune(u.Phone)) != 13 {
		return false
	}
	if len([]rune(u.Mail)) < 7 || len([]rune(u.Mail)) > 40 {
		return false
	}
	if len([]rune(u.Fullname)) < 7 || len([]rune(u.Fullname)) > 40 {
		return false
	}
	return true
}
