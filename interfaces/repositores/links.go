package repositories

import (
	"database/sql"
	"errors"
	"fmt"
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

func (l *links) MarkAsUsed(token string) error {
	_, err := l.db.Exec("UPDATE upload_links SET used = TRUE WHERE token = $1", token)

	if err != nil {
		log.Printf("marking link as used: %v", err)
		return err
	}
	return nil
}

// ValidateToken Validate the token and check expiration
func (l *links) ValidateToken(token string) error {
	var expiration time.Time
	err := l.db.Get(&expiration, "SELECT expiration FROM upload_links WHERE token = $1 AND used = FALSE", token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("invalid or expired token: %w", err)
		}
		return fmt.Errorf("error getting token from database: %w", err)
	}

	if time.Now().After(expiration) {
		return fmt.Errorf("token has expired: %w", err)
	}

	return nil
}
