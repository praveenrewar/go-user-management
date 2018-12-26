package jwtauthenticate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
)

//UserJWTData is used to get user_id from jwt
type UserJWTData struct {
	UserID string `json:"user_id"`
}

//Authenticate is used to verify jwt token
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type AuthenticationMessage struct {
			Status  int32  `json:"status"`
			Message string `json:"message"`
		}
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			token := strings.Fields(authorizationHeader)
			if len(token) > 1 {
				JWT := token[1]
				type Exception struct {
					Message string `json:"message"`
				}

				token, _ := jwt.Parse(JWT, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					var user UserJWTData
					mapstructure.Decode(claims, &user)
					context.Set(r, "user_jwt", user)
					next.ServeHTTP(w, r)
				} else {
					w.Header().Set("Content-type", "application/json")
					w.WriteHeader(401)
					message := AuthenticationMessage{
						Status:  401,
						Message: "Invalid authorization token"}
					userJSON, _ := json.Marshal(message)
					w.Write(userJSON)
					return
				}
			} else {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(400)
				message := AuthenticationMessage{
					Status:  400,
					Message: "Please insert authorization header in the format 'Bearer {token}'"}
				userJSON, _ := json.Marshal(message)
				w.Write(userJSON)
				return
			}
		} else {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(400)
			message := AuthenticationMessage{
				Status:  400,
				Message: "An authorization header is required"}
			userJSON, _ := json.Marshal(message)
			w.Write(userJSON)
			return
		}

	})
}
