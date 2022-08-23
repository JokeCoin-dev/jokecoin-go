package mpt

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core/database"
	"jokecoin-go/core/utils"
	_ "net/http/pprof"
	"testing"
)

func TestMPT(t *testing.T) {
	err := database.InitRocksDB("test_rocksdb")
	utils.PanicIfErr(err)
	root := NewTrie()
	for i := 0; i < 10000; i++ {
		key := sha3.Sum256([]byte(fmt.Sprintf("key%d", i)))
		value := []byte(fmt.Sprintf("value%d", i))
		root, err = root.Put(key[:], value)
		utils.PanicIfErr(err)
		v2, err := root.Get(key[:])
		utils.PanicIfErr(err)
		utils.Assert(bytes.Compare(v2, value) == 0)
	}
	for i := 0; i < 10000; i++ {
		key := sha3.Sum256([]byte(fmt.Sprintf("key%d", i)))
		value := []byte(fmt.Sprintf("value%d", i))
		v2, err := root.Get(key[:])
		utils.PanicIfErr(err)
		utils.Assert(bytes.Compare(v2, value) == 0)
	}
	db := database.GetDB()
	db.Close()
}
