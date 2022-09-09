package core

import (
	"jokecoin-go/core/block"
	"jokecoin-go/core/common"
	"jokecoin-go/core/config"
	"jokecoin-go/core/database"
	"jokecoin-go/core/errors"
	"jokecoin-go/core/merkle"
	"jokecoin-go/core/mpt"
	"jokecoin-go/core/utils"
	"log"
	"sort"
	"time"
)

type Chain struct {
	config           config.NodeConfig
	gConfig          config.NodeGlobalConfig
	HighestChain     []*block.Block
	UnresolvedBlocks map[common.Hash]*block.Block
	Son              map[common.Hash][]common.Hash
	TxPool           map[common.Hash]*block.Transaction
}

func NewChain(cfg config.NodeConfig, gcfg config.NodeGlobalConfig) (*Chain, error) {
	log.Println("Creating a new chain...")
	db := database.GetDB()
	cn := &Chain{config: cfg, gConfig: gcfg}
	if !block.CheckGenesisBlock(gcfg.GenesisBlock) {
		log.Fatalln("Invalid genesis block")
	}
	cn.TxPool = make(map[common.Hash]*block.Transaction)
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
	cn.TxPool = make(map[common.Hash]*block.Transaction)
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

func (cn *Chain) InsertBlock(b *block.Block) error {
	if b.Header.Height == 0 {
		// Cannot insert genesis block
		return errors.New(errors.DataConsistencyError)
	}
	if b.Header.ParentHash == cn.HighestChain[len(cn.HighestChain)-1].Header.ComputeHash() {
		// This block is the son of the highest block
		err := b.WriteDB()
		if err != nil {
			return err
		}
		return cn.ExtendHighestChain(b)
	}
	return nil
}

func (cn *Chain) ExtendHighestChain(b *block.Block) error {
	pa := cn.HighestChain[len(cn.HighestChain)-1]
	if b.Header.Height != pa.Header.Height+1 || b.Header.ParentHash != pa.Header.ComputeHash() {
		return errors.New(errors.DataConsistencyError)
	}
	cn.HighestChain = append(cn.HighestChain, b)
	return nil
}

type SortTx struct {
	Txs []*block.Transaction
	P   []int
}

func (s *SortTx) Len() int {
	return len(s.Txs)
}
func (s *SortTx) Swap(i, j int) {
	s.P[i], s.P[j] = s.P[j], s.P[i]
}
func (s *SortTx) Less(i, j int) bool {
	return s.Txs[s.P[i]].Fee > s.Txs[s.P[j]].Fee
}

func (cn *Chain) PackTransactions(lb *block.Block, miner common.Address) ([]block.Transaction, common.Hash) {
	T := SortTx{Txs: make([]*block.Transaction, len(cn.TxPool)), P: make([]int, len(cn.TxPool))}
	i := 0
	for _, tx := range cn.TxPool {
		T.Txs[i] = tx
		T.P[i] = i
		i++
	}
	txs := make([]block.Transaction, 0)
	txs = append(txs, block.CreateCoinbaseTransaction(miner))
	sort.Sort(&T)
	flag := true
	state := mpt.MerklePatriciaTrie{Root: lb.Header.StateHash}
	for flag && len(txs) < block.MAX_TXS_PER_BLOCK {
		flag = false
		for i := 0; i < len(T.Txs); i++ {
			if T.P[i] == -1 {
				continue
			}
			tx := T.Txs[T.P[i]]
			if ns, err := block.ExecuteTransaction(*tx, state); err == nil {
				txs = append(txs, *tx)
				state = ns
				T.P[i] = -1
				flag = true
				break
			}
		}
	}
	return txs, state.Root
}

func (cn *Chain) GetBlockCandidate(miner common.Address) *block.Block {
	lb := cn.HighestChain[len(cn.HighestChain)-1]
	txs, state := cn.PackTransactions(lb, miner)
	bh := block.BlockHeader{
		Height:           lb.Header.Height + 1,
		ParentHash:       lb.Header.ComputeHash(),
		TransactionHash:  merkle.BuildTransactionTree(txs),
		StateHash:        state,
		Time:             time.Now().Unix(),
		ExtraData:        []byte{},
		Difficulty:       lb.Header.Difficulty,
		LastBlockTime:    lb.Header.Time,
		LastKeyBlockTime: lb.Header.LastKeyBlockTime,
	}
	b := &block.Block{Header: bh, TXs: txs}
	return b
}
