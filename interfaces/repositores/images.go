package repositories

import (
	"log"

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
