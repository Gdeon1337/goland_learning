package models

import (
	u "../utils"
	"github.com/jinzhu/gorm"
)

type Currency struct {
	gorm.Model
	Name      string  `json:"name"`
	RubleRate float32 `json:"ruble_rate"`
	Count     float32 `json:"count";gorm:"-"`
}

type CurrencyConvert struct {
	BaseCurrencyId    int     `json:"base_currency_id"`
	ConvertCurrencyId int     `json:"convert_currency_id"`
	Count             float32 `json:"count"`
}

type CurrencyConvertResponse struct {
	BaseCurrency    Currency `json:"base_currency"`
	ConvertCurrency Currency `json:"convert_currency"`
}

func (currency *Currency) Create() *Currency {
	GetDB().Create(currency)
	return currency
}

func GetAllCurrencies(limit, offset int) []Currency {
	currencies := make([]Currency, 4)
	GetDB().Limit(limit).Limit(offset).Find(&currencies)
	return currencies
}

func (currency *Currency) Update() *Currency {
	temp := &Currency{}
	GetDB().Table("currencies").Where("name = ?", currency.Name).First(temp)
	GetDB().Model(&temp).Update(currency)
	return temp
}

func (currency *Currency) Delete() *Currency {
	temp := &Currency{}
	GetDB().Table("currencies").Where("id = ?", int(currency.ID)).First(temp)
	GetDB().Delete(&temp)
	return temp
}

func (currencyConvert *CurrencyConvert) Convert() map[string]interface{} {
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
