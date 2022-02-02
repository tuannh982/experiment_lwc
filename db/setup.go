package db

import (
	"experiment_lwc/commons"
	"github.com/boltdb/bolt"
	"time"
)

type DB struct {
	DB *bolt.DB
}

var dbOpts = &bolt.Options{Timeout: 10 * time.Second}

func Setup(path string) *DB {
	db, err := bolt.Open(path, 0600, dbOpts)
	commons.Must(err)
	commons.Must(db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("containers"))
		commons.Must(err)
		_, err = tx.CreateBucketIfNotExists([]byte("networks"))
		commons.Must(err)
		return nil
	}))
	return &DB{DB: db}
}
