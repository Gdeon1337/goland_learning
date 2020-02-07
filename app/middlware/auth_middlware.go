package middlware

import (
	"../../models"
	u "../../utils"
	"context"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Check Auth Token")
		notAuth := []string{"/api/user/login"}
		checkAdmin := []string{"/api/user", "/admin"}
		requestPath := r.URL.Path

		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string] interface{})
		tokenHeader := r.Header.Get("Authorization")
		tokenCookie, _ := r.Cookie("Authorization")
		if tokenCookie == nil && tokenHeader == "" {
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response, 403)
			return
		}
		if tokenHeader == "" && tokenCookie != nil {
			tokenHeader = tokenCookie.Value
		}
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response, 403)
			return
		}

		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response, 403)
			return
		}

		if !token.Valid {
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response, 403)
			return
		}
		user, ok := models.GetUser(int(tk.UserId))
		if !ok {
			response = u.Message(false, "User is not active.")
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response, http.StatusForbidden)
			return
		}

		for _, value := range checkAdmin {
			if strings.Contains(requestPath, value) && strings.ToLower(user.Role) != "admin" {
				response = u.Message(false, "Permission denied.")
				w.Header().Add("Content-Type", "application/json")
				u.Respond(w, response, http.StatusForbidden)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserId)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}