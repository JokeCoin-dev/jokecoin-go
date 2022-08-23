package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

var ldb *leveldb.DB

func initLevelDB(path string) *leveldb.DB {
	if ldb != nil {
		log.Panicln("LevelDB already initialized")
	}
	var err error
	ldb, err = leveldb.OpenFile(path, nil)
	if err != nil {
		log.Panicln(err)
	}
	return ldb
}
func closeLevelDB() {
	if ldb == nil {
		log.Panicln("DB is not initialized")
		return
	}
	err := ldb.Close()
	if err != nil {
		log.Panicln(err)
	}
}

func getLevelDB() *leveldb.DB {
	return ldb
}
