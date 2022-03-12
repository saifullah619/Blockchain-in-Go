package main

import (
	a1 "github.com/adan7950/assignment01IBC"
)

func main() {
	var chainHead *a1.Block
	genesis := a1.BlockData{Transactions: []string{"S2E", "S2Z"}}
	chainHead = a1.InsertBlock(genesis, chainHead)
	secondBlock := a1.BlockData{Transactions: []string{"E2Alice", "E2Bob", "S2John"}}
	chainHead = a1.InsertBlock(secondBlock, chainHead)

	a1.ListBlocks(chainHead)
	a1.ChangeBlock("S2E", "S2Trudy", chainHead)
	a1.ListBlocks(chainHead)
	a1.VerifyChain(chainHead)
}
