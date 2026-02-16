package api

import (
	"net/http"
	"os"
)

func logging(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs.Print("Request: %s %s", r.Method, r.URL.Path)
		next(w, r)
	})
}
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if len(password) > 0 {
			var jwt string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}
			valid, err := tokenIsValid(jwt)
			if err != nil {
				writeError(w, err, http.StatusBadRequest)
				return
			}
			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
