package app

import (
	"../models"
	"./controllers"
	"./middlware"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func Run(){
	router := mux.NewRouter()
	router.Use(middlware.JwtAuthentication)
	router.HandleFunc("/api/user",
		controllers.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/login",
		controllers.Authenticate).Methods("POST")

	router.HandleFunc("/api/currency",
		controllers.CreateCurrencies).Methods("POST")
	router.HandleFunc("/api/currency",
		controllers.UpdateCurrencies).Methods("PUT")
	router.HandleFunc("/api/currency",
		controllers.DeleteCurrencies).Methods("DELETE")
	router.HandleFunc("/api/currency/convert",
		controllers.ConvertCurrency).Methods("GET")
	router.HandleFunc("/api/currency",
		controllers.GetCurrencies).Methods("GET")

	adminMux := http.NewServeMux()
	models.Admin.MountTo("/admin", adminMux)

	router.PathPrefix("/admin").Handler(adminMux)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	err := http.ListenAndServe(":" + port, router)

	if err != nil {
		fmt.Print(err)
	}
}