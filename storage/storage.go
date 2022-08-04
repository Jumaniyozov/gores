package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage struct {
	config *Config
	db     *sql.DB
}

func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (storage *Storage) Open() error {
	db, err := sql.Open("postgres", storage.config.DatabaseURI)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	storage.db = db

	return nil
}

func (storage *Storage) Close() error {
	if err := storage.db.Close(); err != nil {
		return err
	}
	return nil
}
