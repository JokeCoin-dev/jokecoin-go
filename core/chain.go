package core

import (
	"jokecoin-go/core/block"
	"jokecoin-go/core/common"
	"jokecoin-go/core/config"
	"jokecoin-go/core/database"
	"jokecoin-go/core/errors"
	"jokecoin-go/core/utils"
	"log"
)

type Chain struct {
	config           config.NodeConfig
	gConfig          config.NodeGlobalConfig
	HighestChain     []*block.Block
	UnresolvedBlocks map[common.Hash]*block.Block
	Son              map[common.Hash][]common.Hash
}

func NewChain(cfg config.NodeConfig, gcfg config.NodeGlobalConfig) (*Chain, error) {
	log.Println("Creating a new chain...")
	db := database.GetDB()
	cn := &Chain{config: cfg, gConfig: gcfg}
	if !block.CheckGenesisBlock(gcfg.GenesisBlock) {
		log.Fatalln("Invalid genesis block")
	}
	cn.UnresolvedBlocks = make(map[common.Hash]*block.Block)
	cn.Son = make(map[common.Hash][]common.Hash)
	cn.HighestChain = make([]*block.Block, 0)
	cn.HighestChain = append(cn.HighestChain, gcfg.GenesisBlock)
	err := gcfg.GenesisBlock.WriteDB()
	if err != nil {
		return nil, err
	}
	genesisHash := gcfg.GenesisBlock.Header.ComputeHash()
	err = db.Put([]byte("highest_chain_block"), genesisHash[:])
	if err != nil {
		return nil, errors.Wrap(err, errors.DatabaseError)
	}
	return cn, nil
}

func LoadChain(cfg config.NodeConfig, gcfg config.NodeGlobalConfig) *Chain {
	log.Println("Loading chain...")
	db := database.GetDB()
	cn := &Chain{config: cfg, gConfig: gcfg}
	cn.UnresolvedBlocks = make(map[common.Hash]*block.Block)
	cn.Son = make(map[common.Hash][]common.Hash)
	highestBlockHash := db.MustGet([]byte("highest_chain_block"))
	highestBlock, err := block.LoadBlock(common.ToHash(highestBlockHash))
	utils.PanicIfErr(err)
	cn.HighestChain = make([]*block.Block, highestBlock.Header.Height+1)
	cn.HighestChain[highestBlock.Header.Height] = highestBlock
	for i := highestBlock.Header.Height - 1; i >= 0; i-- {
		b, err := block.LoadBlock(common.ToHash(cn.HighestChain[i+1].Header.ParentHash[:]))
		utils.PanicIfErr(err)
		cn.HighestChain[i] = b
	}
	log.Printf("Chain loaded, highest block height: %d", highestBlock.Header.Height)
	return cn
}
