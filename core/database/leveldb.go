package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

type LevelDB struct {
	db *leveldb.DB
}

func (l *LevelDB) Put(key []byte, value []byte) error {
	return l.db.Put(key, value, nil)
}

func (l *LevelDB) Get(key []byte) ([]byte, error) {
	return l.db.Get(key, nil)
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
