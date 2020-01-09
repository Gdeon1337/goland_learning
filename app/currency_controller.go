package app

import (
	"../models"
	u "../utils"
	"encoding/json"
	"net/http"
	"strconv"
)

var CreateCurrencies = func(w http.ResponseWriter, r *http.Request) {
	currency := &models.Currency{}
	err := json.NewDecoder(r.Body).Decode(currency)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := currency.Create()
	u.Respond(w, resp)
}


var UpdateCurrencies = func(w http.ResponseWriter, r *http.Request) {
	currency := &models.Currency{}
	err := json.NewDecoder(r.Body).Decode(currency)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := currency.Update()
	u.Respond(w, resp)
}

var ConvertCurrency = func(w http.ResponseWriter, r *http.Request) {
	convertCurrency := &models.CurrencyConvert{}
	rawCurrencyConvert, ok := r.URL.Query()["convert_currency_id"]
	if !ok {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}else{
		i64, err := strconv.ParseInt(rawCurrencyConvert[0], 10, 32)
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"))
			return
		}
		convertCurrency.ConvertCurrencyId = int(i64)
	}
	rawCurrencyBase, ok := r.URL.Query()["base_currency_id"]
	if !ok {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}else{
		i64, err := strconv.ParseInt(rawCurrencyBase[0], 10, 32)
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"))
			return
		}
		convertCurrency.BaseCurrencyId = int(i64)
	}
	rawCount, ok := r.URL.Query()["count"]
	if !ok {
		convertCurrency.Count = 1
	}else{
		i64, err := strconv.ParseInt(rawCount[0], 10, 32)
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"))
			return
		}
		convertCurrency.Count = float32(i64)
	}
	resp := convertCurrency.Convert()
	u.Respond(w, resp)
}
