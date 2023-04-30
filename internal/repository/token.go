package repository

import (
	"github.com/boltdb/bolt"
	"telegramBot/internal/config"
)

type Bucket string

const (
	AccessToken  Bucket = "access_token"
	RequestToken Bucket = "request_token"
)

type TokenRepository interface {
	Save(chatID int64, token string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
}

func InitBoltDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(AccessToken))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(RequestToken))
		if err != nil {
			return err
		}
		return nil

	}); err != nil {
		return nil, err
	}

	return db, nil
}
