package handlers

import (
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
	SaveMetadata(filePath string, err error, dimensions string, cameraModel string, location string)
}

type StatisticsRepository interface {
	UpdateUploadStatistics()
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
				go extractAndSaveMetadata(imageRepo, filePath)

				// Update upload statistics
				go statisticsRepo.UpdateUploadStatistics()
			}
		}

		// Mark token as used after successful upload
		if err = linksRepo.MarkAsUsed(token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Files uploaded successfully"))
	}
}

func extractAndSaveMetadata(imageRepo ImageRepository, filePath string) {
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

	imageRepo.SaveMetadata(filePath, err, dimensions, cameraModel, location)
}
