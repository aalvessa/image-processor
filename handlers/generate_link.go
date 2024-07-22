package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

type GenerateLinkRequest struct {
	SecretToken string `json:"secret_token"`
	Expiration  int    `json:"expiration"`
}

type GenerateLinkResponse struct {
	UploadLink string `json:"upload_link"`
}

var jwtKey = []byte("your_secret_key")

func GenerateUploadLink(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GenerateLinkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		expirationTime := time.Now().Add(time.Duration(req.Expiration) * time.Minute)
		claims := &jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO upload_links (token, expiration) VALUES ($1, $2)", tokenString, expirationTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(GenerateLinkResponse{UploadLink: "/upload?token=" + tokenString})
	}
}
