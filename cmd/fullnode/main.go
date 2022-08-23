package main

import (
	"fmt"
	"golang.org/x/crypto/sha3"
	"jokecoin-go/core/block"
	"jokecoin-go/core/database"
)

func main() {
	fmt.Println("Hello World!")
	fmt.Println(sha3.NewLegacyKeccak256().Size())
	a := sha3.NewLegacyKeccak256().Sum([]byte(""))
	for i := range a {
		fmt.Printf("%02x", a[i])
	}
	h := block.BlockHeader{}
	h.ComputeHash()
	database.InitLevelDB("db")
}
