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
		u.Respond(w, u.Message(false, "Invalid request"), 415)
		return
	}
	resp := currency.Create()
	u.Respond(w, resp, 200)
}

var GetCurrencies = func(w http.ResponseWriter, r *http.Request) {
	rawLimit := r.URL.Query().Get("limit")
	rawOffset := r.URL.Query().Get("offset")
	limit := 1
	offset := 10
	if rawLimit != "" {
		rawLimit, errParseLimit := strconv.ParseInt(rawLimit, 10, 32)
		if errParseLimit != nil {
			u.Respond(w, u.Message(false, "Invalid request"), 415)
			return
		}
		limit = int(rawLimit)
	}
	if rawOffset != "" {
		rawOffset, errParseOffset := strconv.ParseInt(rawOffset, 10, 32)
		if errParseOffset != nil {
			u.Respond(w, u.Message(false, "Invalid request"), 415)
			return
		}
		offset = int(rawOffset)
	}
	resp := models.GetAllCurrencies(limit, offset)
	u.Respond(w, resp, 200)
}


var UpdateCurrencies = func(w http.ResponseWriter, r *http.Request) {
	currency := &models.Currency{}
	err := json.NewDecoder(r.Body).Decode(currency)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"), 415)
		return
	}
	resp := currency.Update()
	u.Respond(w, resp, 200)
}

var DeleteCurrencies = func(w http.ResponseWriter, r *http.Request) {
	rawCurrencyId := r.URL.Query().Get("currency_id")
	currency := &models.Currency{}
	if rawCurrencyId == "" {
		u.Respond(w, u.Message(false, "Invalid request"), 404)
		return
	}else{
		i64, err := strconv.ParseInt(rawCurrencyId, 10, 32)
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"), 415)
			return
		}
		currency.ID = uint(i64)
	}
	resp := currency.Delete()
	u.Respond(w, resp, 200)
}

var ConvertCurrency = func(w http.ResponseWriter, r *http.Request) {
	convertCurrency := &models.CurrencyConvert{}
	rawCurrencyConvert := r.URL.Query().Get("convert_currency_id")
	rawCurrencyBase := r.URL.Query().Get("base_currency_id")
	if rawCurrencyConvert == "" || rawCurrencyBase == "" {
		u.Respond(w, u.Message(false, "Invalid request"), 415)
		return
	}else{
		CurrencyConvert, errParseConvert := strconv.ParseInt(rawCurrencyConvert, 10, 32)
		CurrencyBase, errParseBase := strconv.ParseInt(rawCurrencyBase, 10, 32)
		if errParseConvert != nil || errParseBase != nil {
			u.Respond(w, u.Message(false, "Invalid request"), 415)
			return
		}
		convertCurrency.BaseCurrencyId = int(CurrencyBase)
		convertCurrency.ConvertCurrencyId = int(CurrencyConvert)
	}
	rawCount := r.URL.Query().Get("count")
	if rawCount == "" {
		convertCurrency.Count = 1
	}else{
		i64, err := strconv.ParseInt(rawCount, 10, 32)
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"), 415)
			return
		}
		convertCurrency.Count = float32(i64)
	}
	resp := convertCurrency.Convert()
	u.Respond(w, resp, 200)
}
