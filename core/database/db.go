package database

import "log"

type DB interface {
	Put(key []byte, value []byte) error
	MustPut(key []byte, value []byte)
	// Get return (nil,nil) if key not found
	Get(key []byte) ([]byte, error)
	MustGet(key []byte) []byte
	Close() error
}

var db DB

func GetDB() DB {
	if db == nil {
		log.Panic("DB is not initialized")
	}
	return db
}
