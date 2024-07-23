package repositories

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type links struct {
	db *sqlx.DB
}

func NewLinks(db *sqlx.DB) *links {
	return &links{
		db: db,
	}
}

func (l *links) CreateLink(token string, expiration time.Time) error {
	_, err := l.db.Exec("INSERT INTO upload_links (token, expiration) VALUES ($1, $2)", token, expiration)

	if err != nil {
		log.Printf("creating link: %v", err)
		return err
	}

	return nil
}
