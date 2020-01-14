package models

import (
	u "../utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"token";gorm:"-"`
	Role     string `json:"role"`
}

func (user *User) Validate() (map[string]interface{}, bool) {

	if len(user.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	temp := &User{}

	err := GetDB().Table("users").Where("login = ?", user.Login).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Login != "" {
		return u.Message(false, "Login address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}


func GetUser (UserId int) (*User, bool) {
	temp := &User{}
	err := GetDB().Table("users").Where("id = ?", UserId).First(temp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return temp , false
		}
		return temp, false
	}
	return temp, true
}


func (user *User) Create() map[string]interface{} {

	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	GetDB().Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create user, connection error.")
	}

	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = ""

	response := u.Message(true, "user has been created")
	response["user"] = user
	return response
}

func Login(login, password string) map[string]interface{} {

	user := &User{}
	err := GetDB().Table("users").Where("login = ?", login).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Login not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	user.Password = ""

	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString // Сохраните токен в ответе

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}
