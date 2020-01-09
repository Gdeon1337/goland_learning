package main

import (
	"./app"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

func task() {
	timerCh := time.Tick(time.Duration(60) * time.Second)
	for range timerCh {
		app.CurrencyParser()
	}
}


func main() {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)
	go task()

	router.HandleFunc("/api/user/new",
		app.CreateUser).Methods("POST")

	router.HandleFunc("/api/user/login",
		app.Authenticate).Methods("POST")

	router.HandleFunc("/api/currency/new",
		app.CreateCurrencies).Methods("POST")

	router.HandleFunc("/api/currency/update",
		app.UpdateCurrencies).Methods("POST")

	router.HandleFunc("/api/currency/convert",
		app.ConvertCurrency).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	err := http.ListenAndServe(":" + port, router)

	if err != nil {
		fmt.Print(err)
	}
}
