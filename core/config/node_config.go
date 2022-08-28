package config

import "jokecoin-go/core/block"

type NodeConfig struct {
	Database     string `json:"database"`
	DatabasePath string `json:"database_path"`
}

type NodeGlobalConfig struct {
	ChainID        int64        `json:"chain_id"`
	GenesisBlock   *block.Block `json:"genesis_block"`
	BootstrapNodes []string     `json:"bootstrap_nodes"`
}
