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

func (i *images) SaveMetadata(
	filePath string,
	dimensions string,
	cameraModel string,
	location string,
) (int64, error) {
	// Save metadata to the database
	var imageID int64
	err := i.db.QueryRow(`
        INSERT INTO images (path, dimensions, camera_model, location, format)
        VALUES ($1, $2, $3, $4, $5) RETURNING id
    `, filePath, dimensions, cameraModel, location, filepath.Ext(filePath)).Scan(&imageID)
	if err != nil {
		return 0, fmt.Errorf("error saving metadata to database: %v", err)
	}

	return imageID, nil
}
