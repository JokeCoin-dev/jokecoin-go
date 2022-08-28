package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"jokecoin-go/core/errors"
	"log"
)

type LevelDB struct {
	db *leveldb.DB
}

func (l *LevelDB) Put(key []byte, value []byte) error {
	return l.db.Put(key, value, nil)
}

func (l *LevelDB) MustPut(key []byte, value []byte) {
	err := l.Put(key, value)
	if err != nil {
		log.Panicf("Database error: %v\n", err)
	}
}

func (l *LevelDB) Get(key []byte) ([]byte, error) {
	v, err := l.db.Get(key, nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		return nil, nil
	}
	return v, err
}

func (l *LevelDB) MustGet(key []byte) []byte {
	value, err := l.Get(key)
	if err != nil {
		log.Panicf("Database error: %v\n", err)
	}
	return value
}

func (l *LevelDB) Close() error {
	return l.db.Close()
}

func InitLevelDB(path string) error {
	if db != nil {
		log.Panicln("Database already initialized")
	}
	var err error
	ldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return err
	}
	db = &LevelDB{db: ldb}
	return nil
}
