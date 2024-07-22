package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"image/jpeg"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

type UploadImageResponse struct {
	IDs []int `json:"ids"`
}

func UploadImage(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			http.Error(w, "token is required", http.StatusBadRequest)
			return
		}

		var uploadLink struct {
			Used bool
		}
		err := db.Get(&uploadLink, "SELECT used FROM upload_links WHERE token=$1", tokenString)
		if err != nil || uploadLink.Used {
			http.Error(w, "invalid or used token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		r.ParseMultipartForm(10 << 20) // 10 MB limit
		form := r.MultipartForm
		files := form.File["images"]
		var ids []int

		for _, file := range files {
			filePath := "uploads/" + filepath.Base(file.Filename)
			src, err := file.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer src.Close()

			f, err := os.Create(filePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			_, err = src.Seek(0, 0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			img, err := jpeg.Decode(src)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			exifData, _ := exif.Decode(src)

			cameraModel, _ := exifData.Get(exif.Model)
			//lat, lng, _ := exifData.LatLong()

			newPath := "uploads/resized_" + filepath.Base(file.Filename)
			newFile, err := os.Create(newPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer newFile.Close()

			resizedImage := resize.Resize(800, 0, img, resize.Lanczos3)
			jpeg.Encode(newFile, resizedImage, nil)

			var imageID int
			err = db.QueryRow("INSERT INTO images (path, dimensions, camera_model, location) VALUES ($1, $2, $3, $4) RETURNING id",
				newPath, "800x600", cameraModel, "latLng").Scan(&imageID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ids = append(ids, imageID)
		}

		db.Exec("UPDATE upload_links SET used = true WHERE token = $1", tokenString)
		json.NewEncoder(w).Encode(UploadImageResponse{IDs: ids})
	}
}
