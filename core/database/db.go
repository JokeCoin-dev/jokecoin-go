package database

import "log"

type DB interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Close() error
}

var db DB

func GetDB() DB {
	if db == nil {
		log.Panic("DB is not initialized")
	}
	return db
}
