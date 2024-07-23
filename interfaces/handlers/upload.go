package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func UploadImage(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Token is required", http.StatusBadRequest)
			return
		}

		// Validate the token and check expiration
		var expiration time.Time
		err = db.Get(&expiration, "SELECT expiration FROM upload_links WHERE token = $1 AND used = FALSE", token)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if time.Now().After(expiration) {
			http.Error(w, "Token has expired", http.StatusUnauthorized)
			return
		}

		// Ensure the uploads directory exists
		uploadDir := "uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err := os.Mkdir(uploadDir, os.ModePerm)
			if err != nil {
				http.Error(w, fmt.Sprintf("could not create upload directory: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// Process each uploaded file
		for _, headers := range r.MultipartForm.File {
			for _, header := range headers {
				file, err := header.Open()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer file.Close()

				filePath := filepath.Join(uploadDir, header.Filename)
				outFile, err := os.Create(filePath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer outFile.Close()

				_, err = io.Copy(outFile, file)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// Extract and save image metadata
				go extractAndSaveMetadata(db, filePath)

				// Update upload statistics
				go updateUploadStatistics(db)
			}
		}

		// Mark token as used after successful upload
		_, err = db.Exec("UPDATE upload_links SET used = TRUE WHERE token = $1", token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Files uploaded successfully"))
	}
}

func extractAndSaveMetadata(db *sqlx.DB, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	defer f.Close()

	exif.RegisterParsers(mknote.All...)

	var cameraModel, dimensions, location string

	x, err := exif.Decode(f)
	if err != nil {
		fmt.Printf("error decoding EXIF data: %v\n", err)
	} else {
		// Extract metadata fields
		model, _ := x.Get(exif.Model)
		xDimension, _ := x.Get(exif.PixelXDimension)
		yDimension, _ := x.Get(exif.PixelYDimension)
		cameraModel = model.String()
		dimensions = fmt.Sprintf("%dx%d", xDimension.Count, yDimension.Count)
		location = "" // extract location
	}

	// Save metadata to the database
	_, err = db.Exec(`
        INSERT INTO images (path, dimensions, camera_model, location, format)
        VALUES ($1, $2, $3, $4, $5)
    `, filePath, dimensions, cameraModel, location, filepath.Ext(filePath))
	if err != nil {
		fmt.Printf("error saving metadata to database: %v\n", err)
	}
}

func updateUploadStatistics(db *sqlx.DB) {
	currentDate := time.Now().Format("2006-01-02")
	_, err := db.Exec(`
        INSERT INTO upload_statistics (upload_date, upload_count)
        VALUES ($1, 1)
        ON CONFLICT (upload_date)
        DO UPDATE SET upload_count = upload_statistics.upload_count + 1
    `, currentDate)
	if err != nil {
		fmt.Printf("error updating upload statistics: %v\n", err)
	}
}
