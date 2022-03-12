package assignment01IBC

import (
	"crypto/sha256"
	"fmt"
)

type BlockData struct {
	Transactions []string
}
type Block struct {
	Data        BlockData
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

func CalculateHash(inputBlock *Block) string {
	//Calculate Block Hash by the internal Data
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v", inputBlock))))
}
func InsertBlock(dataToInsert BlockData, chainHead *Block) *Block {

	if chainHead == nil {
		// Create new Chainhead
		chainHead = new(Block)
		chainHead.PrevPointer = nil
		chainHead.PrevHash = ""
		chainHead.Data = dataToInsert
		chainHead.CurrentHash = CalculateHash(chainHead)
		return chainHead
	}
	var newBlock = new(Block)
	newBlock.PrevPointer = chainHead
	newBlock.Data = dataToInsert
	newBlock.PrevHash = chainHead.CurrentHash
	newBlock.CurrentHash = CalculateHash(newBlock)
	return newBlock

}

func ChangeBlock(oldTrans string, newTrans string, chainHead *Block) {
	for chainHead != nil {
		for index, value := range chainHead.Data.Transactions {
			if value == oldTrans {
				chainHead.CurrentHash = ""
				chainHead.Data.Transactions[index] = newTrans
				chainHead.CurrentHash = CalculateHash(chainHead)
			}
		}
		chainHead = chainHead.PrevPointer
	}
}
func ListBlocks(chainHead *Block) {
	it := 1
	for chainHead != nil {
		fmt.Println(it, ": ", chainHead.Data.Transactions)
		it += 1
		chainHead = chainHead.PrevPointer
	}
	fmt.Println("")
}
func VerifyChain(chainHead *Block) {
	for chainHead.PrevPointer != nil {
		if chainHead.PrevPointer.CurrentHash != chainHead.PrevHash {
			fmt.Println("Block Chain Compromised")
			return
		}
		chainHead = chainHead.PrevPointer
	}
	fmt.Println("Block Chain Verified")
}
