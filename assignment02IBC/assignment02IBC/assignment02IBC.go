package assignment02IBC

import (
	"crypto/sha256"
	"fmt"
)

const miningReward = 100
const rootUser = "Satoshi"

type BlockData struct {
	Title    string
	Sender   string
	Receiver string
	Amount   int
}
type Block struct {
	Data        []BlockData
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

func CalculateBalance(userName string, chainHead *Block) int {
	// Get the balance of a person in the blockChain
	balance := 0
	for chainHead != nil {
		for _, data := range chainHead.Data {
			if data.Sender == userName {
				balance -= data.Amount
			}
			if data.Receiver == userName {
				balance += data.Amount
			}
		}
		chainHead = chainHead.PrevPointer
	}
	return balance
}
func CalculateHash(inputBlock *Block) string {
	var blockString = fmt.Sprintf("%v", inputBlock)
	var hash = fmt.Sprintf("%x", sha256.Sum256([]byte(blockString)))
	return hash
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
func VerifyTransaction(transaction *BlockData, chainHead *Block) bool {
	currentBalance := CalculateBalance(transaction.Sender, chainHead)
	if transaction.Sender != "System" && transaction.Amount > currentBalance {
		// Complete transaction in Invalid. No new data can be inserted
		fmt.Printf("Error: %v has %d coins - %d were needed!\n", transaction.Sender, currentBalance, transaction.Amount)
		return false
	}
	return true
}
func InsertBlock(blockData []BlockData, chainHead *Block) *Block {
	// Iterate over all the blockdata's and check for any invalid transaction
	for _, transaction := range blockData {

		if !VerifyTransaction(&transaction, chainHead) {
			return chainHead
		}

		currentBalance := CalculateBalance(transaction.Sender, chainHead)
		//also calculate the balance for current block
		for _, tempITR := range blockData {
			if tempITR.Sender == transaction.Sender && transaction != tempITR {
				currentBalance -= tempITR.Amount
			}
		}

		if transaction.Sender != "System" && transaction.Amount > currentBalance {
			// Complete transaction in Invalid. No new data can be inserted
			fmt.Printf("Error: %v has %d coins - %d were needed!\n", transaction.Sender, currentBalance, transaction.Amount)
			return chainHead
		}
	}

	blockData = append(blockData, BlockData{Title: "Coinbase", Sender: "System", Receiver: rootUser, Amount: miningReward})
	if chainHead == nil {
		// Create new Chainhead
		chainHead = new(Block)
		chainHead.PrevPointer = nil
		chainHead.PrevHash = ""
		chainHead.Data = blockData
		chainHead.CurrentHash = CalculateHash(chainHead)
		return chainHead
	} else {
		var newBlock = new(Block)
		newBlock.PrevPointer = chainHead
		newBlock.Data = blockData
		newBlock.PrevHash = chainHead.CurrentHash
		newBlock.CurrentHash = CalculateHash(newBlock)
		return newBlock
	}
}

func ListBlocks(chainHead *Block) {
	fmt.Println("\n              			     Head")
	fmt.Println("				       ↓		")
	for chainHead != nil {
		for _, block := range chainHead.Data {
			fmt.Printf("Title: %-15v Sender: %-10v Reciever: %-10v Amount: %d", block.Title, block.Sender, block.Receiver, block.Amount)
			fmt.Println()
		}
		chainHead = chainHead.PrevPointer
		if chainHead != nil {
			fmt.Println("				       ↓		")
			fmt.Println()
		}
	}
	fmt.Println()
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
func PremineChain(chainHead *Block, numBlocks int) *Block {
	if chainHead != nil {
		fmt.Println("Error: Can't Premine. Head already exists")
	}
	for i := 0; i < numBlocks; i++ {
		newBlock := []BlockData{{Title: "Premined", Sender: "nil", Receiver: "nil", Amount: 0}}
		chainHead = InsertBlock(newBlock, chainHead)
	}
	return chainHead
}
