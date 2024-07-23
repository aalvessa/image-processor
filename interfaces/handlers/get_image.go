package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ImageGetter interface {
	GetPath(imageID string) (string, error)
}

func GetImage(getter ImageGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		imageID := chi.URLParam(r, "id")

		var imagePath string
		//err := db.Get(&imagePath, "SELECT path FROM images WHERE id = $1", imageID)
		imagePath, err := getter.GetPath(imageID)

		if err != nil {
			http.Error(w, "image not found", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, imagePath)
	}
}
