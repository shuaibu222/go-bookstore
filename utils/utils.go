package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shuaibu222/go-bookstore/auth"
	"github.com/shuaibu222/go-bookstore/config"
)

var jwtSecretKey []byte

func ParseBody(r *http.Request, x interface{}) {
	if body, err := io.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}

func JwtUserIdUsername(w http.ResponseWriter, r *http.Request) (string, string) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Println("Error while loading envs: ", err)
	}

	jwtSecretKey = []byte(config.JWTSecret)
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("No cookie found", err)
		}
		log.Println("Error while getting cookie", err)
	}
	tknStr := c.Value
	claims := &auth.Claims{}

	// token validation
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("jwt signature is invalid", err)
		}
		log.Println("jwt issue", err)
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// id and username from jwt claims
	id := claims.UserID
	userName := claims.Username

	return id, userName
}
