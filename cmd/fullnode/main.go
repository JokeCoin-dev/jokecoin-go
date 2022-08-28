package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core"
	"jokecoin-go/core/block"
	"jokecoin-go/core/config"
	"jokecoin-go/core/database"
	"jokecoin-go/core/utils"
	"log"
	"os"
)

func main() {
	configPath := flag.String("config", "config.json", "config file")
	gConfigPath := flag.String("global-config", "global_config.json", "global config file")
	flag.Parse()
	var cfg config.NodeConfig
	var gcfg config.NodeGlobalConfig
	cf, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v\n", err)
	}
	err = json.Unmarshal(cf, &cfg)
	if err != nil {
		log.Fatalf("Failed to parse config: %v\n", err)
	}
	gcf, err := os.ReadFile(*gConfigPath)
	if err != nil {
		log.Fatalf("Failed to read global config: %v\n", err)
	}
	err = json.Unmarshal(gcf, &gcfg)
	if err != nil {
		log.Fatalf("Failed to parse global config: %v\n", err)
	}
	if cfg.Database == "leveldb" {
		err := database.InitLevelDB(cfg.DatabasePath)
		utils.PanicIfErr(err)
	} else if cfg.Database == "rocksdb" {
		err := database.InitRocksDB(cfg.DatabasePath)
		utils.PanicIfErr(err)
	} else {
		log.Fatalf("Unknown database: %s\n", cfg.Database)
	}
	db := database.GetDB()
	chainIDb := db.MustGet([]byte("current_chain_id"))
	isNewChain := false
	if chainIDb == nil {
		isNewChain = true
	} else {
		var chainID int64
		utils.MustDeserialize(chainIDb, &chainID)
		if chainID != gcfg.ChainID {
			isNewChain = true
		}
	}
	if isNewChain {
		db.MustPut([]byte("current_chain_id"), utils.MustSerialize(gcfg.ChainID))
		_, err = core.NewChain(cfg, gcfg)
		if err != nil {
			log.Fatalf("Failed to initial chain : %v\n", err)
		}
	} else {
		core.LoadChain(cfg, gcfg)
		if err != nil {
			log.Fatalf("Failed to load chain : %v\n", err)
		}
	}
	fmt.Println("Hello World!")
	fmt.Println(sha3.NewLegacyKeccak256().Size())
	a := sha3.NewLegacyKeccak256().Sum([]byte(""))
	for i := range a {
		fmt.Printf("%02x", a[i])
	}
	h := block.BlockHeader{}
	h.ComputeHash()
	database.GetDB().Close()
}
