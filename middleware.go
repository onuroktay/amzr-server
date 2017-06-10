package main

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"github.com/onuroktay/amazon-reader/amzr-server/util"
)

func checkAuth(minRoleValue int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string

		// Read Cookie
		cookie, err := r.Cookie("onurAuth")
		if err != nil {
			sendUnthorized(w, r)
			return
		}

		// Read Token in Cookie -> idUser
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Token invalid")
			}
			return util.GetEncryptionKey(), nil
		})

		// Check if user is authorized
		if err != nil {
			sendUnthorized(w, r)
			return
		}

		// Read IDUser
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID = claims["IdUser"].(string)
		}
		if userID == "" {
			sendUnthorized(w, r)
			return
		}

		// Check if user is authorized
		user, err := database.accesser.GetAccountByIDInDB(userID)
		if err != nil {
			sendUnthorized(w, r)
			return
		}

		if user.RoleValue < minRoleValue {
			sendUnthorized(w, r)
			return
		}

		// Function passed in parameter is executed
		next.ServeHTTP(w, r)
	})
}

func sendUnthorized(w http.ResponseWriter, r *http.Request, ) {
	setHeader(w, r)

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Sorry, you are not authorized"))
}
