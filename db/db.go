package db

import "github.com/boltdb/bolt"

type DB struct {
	db *bolt.DB
}

func (d *DB) Shutdown() error {
	return d.db.Close()
}

func 