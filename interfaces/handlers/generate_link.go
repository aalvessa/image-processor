package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type LinkCreator interface {
	CreateLink(token string, expiration time.Time) error
}

type GenerateLinkRequest struct {
	SecretToken string `json:"secret_token"`
	Expiration  int    `json:"expiration"`
}

type GenerateLinkResponse struct {
	UploadLink string `json:"upload_link"`
}

var jwtKey = []byte("07t2me784edA44fe6a7Ffabf0D534247")

func GenerateUploadLink(creator LinkCreator) http.HandlerFunc {
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

		err = creator.CreateLink(tokenString, expirationTime)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(GenerateLinkResponse{UploadLink: "/upload?token=" + tokenString})
	}
}
