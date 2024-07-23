package repositories

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

type images struct {
	db *sqlx.DB
}

func NewImages(db *sqlx.DB) *images {
	return &images{
		db: db,
	}
}

func (i *images) GetPath(id string) (string, error) {
	var path string
	err := i.db.Get(&path, "SELECT path FROM images WHERE id = $1", id)
	if err != nil {
		log.Println("getting image path:", err)
		return "", err
	}

	return path, nil
}

func (i *images) SaveMetadata(filePath string, err error, dimensions string, cameraModel string, location string) {
	// Save metadata to the database
	_, err = i.db.Exec(`
        INSERT INTO images (path, dimensions, camera_model, location, format)
        VALUES ($1, $2, $3, $4, $5)
    `, filePath, dimensions, cameraModel, location, filepath.Ext(filePath))
	if err != nil {
		fmt.Printf("error saving metadata to database: %v\n", err)
	}
}
