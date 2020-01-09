package app

import (
	"../models"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
)

func CurrencyParser()  {
	currencies := make([]models.Currency, 4)

	models.GetDB().Find(&currencies)
	resp, err := http.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	currenciesParse := gjson.Parse(string(body)).Get("Valute")
	for _, currency := range currencies {
		rubleRate := currenciesParse.Get(currency.Name).Get("Value")
		models.GetDB().Model(&currency).Update("ruble_rate", rubleRate.Float())
	}
}