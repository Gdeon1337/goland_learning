package models

import (
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


func (user *User) Create() *User {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	GetDB().Create(user)
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = ""


	return user
}

func Login(login, password string) (*User, bool) {

	user := &User{}
	err := GetDB().Table("users").Where("login = ?", login).First(user).Error
	if err != nil {
		return nil, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, false
	}
	user.Password = ""

	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString
	return user, true
}