package models

import (
	u "../utils"
	"github.com/jinzhu/gorm"
)


func (currency *Currency) Validate() (map[string]interface{}, bool) {

	temp := &Currency{}

	err := GetDB().Table("currencies").Where("name = ?", currency.Name).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Name != "" {
		return u.Message(false, "Name address already in use."), false
	}

	return u.Message(false, "Requirement passed"), true
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