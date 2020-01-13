package main

import (
	"./app"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func task() {
	timerCh := time.Tick(time.Duration(60) * time.Second)
	for range timerCh {
		log.Print("Start task - update currencies")
		app.CurrencyParser()
	}
}


func main() {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)
	go task()

	router.HandleFunc("/api/user",
		app.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/login",
		app.Authenticate).Methods("POST")

	router.HandleFunc("/api/currency",
		app.CreateCurrencies).Methods("POST")
	router.HandleFunc("/api/currency",
		app.UpdateCurrencies).Methods("PUT")
	router.HandleFunc("/api/currency",
		app.DeleteCurrencies).Methods("DELETE")
	router.HandleFunc("/api/currency/convert",
		app.ConvertCurrency).Methods("GET")
	router.HandleFunc("/api/currency",
		app.GetCurrencies).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	err := http.ListenAndServe(":" + port, router)

	if err != nil {
		fmt.Print(err)
	}
}
