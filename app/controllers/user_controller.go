package controllers

import (
	"../../models"
	u "../../utils"
	"encoding/json"
	"log"
	"net/http"
)


var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	log.Print("CreateUser")
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"), 415)
		return
	}
	if resp, ok := user.Validate(); !ok {
		u.Respond(w, resp, 415)
	}
	response := u.Message(true, "user has been created")
	response["user"] = user.Create()
	log.Print("User user has been created")
	u.Respond(w, resp, 200)
}


var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	log.Print("Auth user")
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"), 415)
		return
	}
	if authUser, ok := models.Login(user.Login, user.Password); !ok {
		resp := u.Message(false, "Invalid login credentials. Please try again")
		u.Respond(w, resp, 403)
	}else{
		resp := u.Message(true, "Logged In")
		resp["user"] = authUser
		u.Respond(w, resp, 200)
	}
}

