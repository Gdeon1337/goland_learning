package models

import (
	u "../utils"
	"github.com/jinzhu/gorm"
)


type Currency struct {
	gorm.Model
	Name      string  `json:"name"`
	RubleRate float32 `json:"ruble_rate"`
	Count float32 `json:"count" gorm:"-"`
}

type CurrencyConvert struct {
	BaseCurrencyId      int  `json:"base_currency_id"`
	ConvertCurrencyId int `json:"convert_currency_id"`
	Count float32 `json:"count"`
}

type CurrencyConvertResponse struct {
	BaseCurrency Currency `json:"base_currency"`
	ConvertCurrency Currency `json:"convert_currency"`
}

func (currency *Currency) Validate() (map[string] interface{}, bool) {

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


func (currency *Currency) Create() map[string] interface{} {

	if resp, ok := currency.Validate(); !ok {
		return resp
	}

	GetDB().Create(currency)
	response := u.Message(true, "currencies has been created")
	response["currency"] = currency
	return response
}

func GetAllCurrencies(limit, offset int) map[string] interface{} {
	currencies := make([]Currency, 4)
	GetDB().Limit(limit).Limit(offset).Find(&currencies)
	response := u.Message(true, "get all currencies")
	response["currencies"] = currencies
	return response
}



func (currency *Currency) Update() map[string] interface{} {
	temp := &Currency{}
	err := GetDB().Table("currencies").Where("name = ?", currency.Name).First(temp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Currency not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}
	GetDB().Model(&temp).Update(currency)
	response := u.Message(true, "currencies has been created")
	response["currency"] = temp
	return response
}


func (currency *Currency) Delete() map[string] interface{} {
	temp := &Currency{}
	err := GetDB().Table("currencies").Where("id = ?", int(currency.ID)).First(temp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Currency not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}
	GetDB().Delete(&temp)
	response := u.Message(true, "currencies has been deleted")
	response["currency"] = temp
	return response
}


func (currencyConvert *CurrencyConvert) Convert() map[string] interface{} {
	baseCurrency := &Currency{}
	convertCurrency := &Currency{}
	err := GetDB().Table("currencies").Where("id = ?", currencyConvert.BaseCurrencyId).First(baseCurrency).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Currency not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}
	baseCurrency.Count = currencyConvert.Count
	err = GetDB().Table("currencies").Where("id = ?", currencyConvert.ConvertCurrencyId).First(convertCurrency).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Currency not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}
	currencyRate := (baseCurrency.RubleRate / convertCurrency.RubleRate) * currencyConvert.Count
	convertCurrency.Count = currencyRate
	response := u.Message(true, "currencies has been converted")
	responseConvertCurrency := &CurrencyConvertResponse{}
	responseConvertCurrency.BaseCurrency = *baseCurrency
	responseConvertCurrency.ConvertCurrency = *convertCurrency
	response["convert"] = responseConvertCurrency
	return response
}