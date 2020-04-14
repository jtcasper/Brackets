package db

import (
	"database/sql"
	"sync"
)

// A DB that can be locked, as SQLite can't be concurrently written to.
type DB struct {
	Db *sql.DB
	Mu *sync.Mutex
}

func New() (*DB, error) {
	db, err := sql.Open("sqlite3", "brackets.sqlite")
	if err != nil {
		return nil, err
	}
	return &DB{Db: db, Mu: &sync.Mutex{}}, err
}
