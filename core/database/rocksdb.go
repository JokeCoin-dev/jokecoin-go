package database

import (
	"github.com/linxGnu/grocksdb"
	"log"
)

type RocksDB struct {
	db *grocksdb.DB
}

var ro *grocksdb.ReadOptions
var wo *grocksdb.WriteOptions

func (r *RocksDB) Put(key []byte, value []byte) error {
	return r.db.Put(wo, key, value)
}

func (r *RocksDB) MustPut(key []byte, value []byte) {
	err := r.Put(key, value)
	if err != nil {
		log.Panicf("Database error: %v\n", err)
	}
}

func (r *RocksDB) Get(key []byte) ([]byte, error) {
	v, err := r.db.Get(ro, key)
	defer v.Free()
	if err != nil {
		return nil, err
	}
	if !v.Exists() {
		return nil, nil
	}
	t := v.Data()
	res := make([]byte, len(t))
	copy(res, t)
	return res, nil
}

func (r *RocksDB) MustGet(key []byte) []byte {
	value, err := r.Get(key)
	if err != nil {
		log.Panicf("Database error: %v\n", err)
	}
	return value
}

func (r *RocksDB) Close() error {
	r.db.Close()
	return nil
}

func InitRocksDB(path string) error {
	if db != nil {
		log.Panicln("Database already initialized")
	}
	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(grocksdb.NewLRUCache(3 << 30))

	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)

	var err error
	rocksDB, err := grocksdb.OpenDb(opts, path)
	if err != nil {
		return err
	}
	db = &RocksDB{db: rocksDB}
	ro = grocksdb.NewDefaultReadOptions()
	wo = grocksdb.NewDefaultWriteOptions()
	return nil
}
