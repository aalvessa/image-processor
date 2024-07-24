package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

type LinksRepository interface {
	MarkAsUsed(token string) error
	ValidateToken(token string) error
}

type ImageRepository interface {
	SaveMetadata(filePath string, dimensions string, cameraModel string, location string) (int64, error)
}

type StatisticsRepository interface {
	UpdateUploadStatistics()
}

type UploadResponse struct {
	ImageIDs []int64 `json:"image_ids"`
}

func UploadImage(linksRepo LinksRepository, imageRepo ImageRepository, statisticsRepo StatisticsRepository) http.HandlerFunc {
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

		if err = linksRepo.ValidateToken(token); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
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

		var imageIDs []int64
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
				imageID, err := extractAndSaveMetadata(imageRepo, filePath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				imageIDs = append(imageIDs, imageID)

				// Update upload statistics
				go statisticsRepo.UpdateUploadStatistics()
			}
		}

		// Mark token as used after successful upload
		if err = linksRepo.MarkAsUsed(token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := UploadResponse{ImageIDs: imageIDs}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func extractAndSaveMetadata(imageRepo ImageRepository, filePath string) (int64, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("error opening file: %v", err)
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

	return imageRepo.SaveMetadata(filePath, dimensions, cameraModel, location)
}
