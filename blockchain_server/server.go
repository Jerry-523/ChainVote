package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Blockchain struct {
	chain               []Block
	currentTransactions []Transaction
	mu                  sync.Mutex
}

type Block struct {
	Index        int           `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int           `json:"proof"`
	PreviousHash string        `json:"previous_hash"`
}

type Transaction struct {
	VoterID     string `json:"voter_id"`
	CandidateID string `json:"candidate_id"`
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		chain:               make([]Block, 0),
		currentTransactions: make([]Transaction, 0),
	}
	bc.NewBlock(100, "1")
	return bc
}

func (bc *Blockchain) GetBlocks() []Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return bc.chain
}

func (bc *Blockchain) NewBlock(proof int, previousHash string) Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	block := Block{
		Index:        len(bc.chain) + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: bc.currentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}

	bc.currentTransactions = []Transaction{}
	bc.chain = append(bc.chain, block)
	return block
}

func (bc *Blockchain) NewTransaction(voterID, candidateID string) int {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.currentTransactions = append(bc.currentTransactions, Transaction{
		VoterID:     voterID,
		CandidateID: candidateID,
	})

	return bc.LastBlock().Index + 1
}

func (bc *Blockchain) LastBlock() Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) ProofOfWork(lastProof int) int {
	proof := 0
	for !validProof(lastProof, proof) {
		proof++
	}
	return proof
}

func Hash(block Block) string {
	blockBytes, _ := json.Marshal(block)
	hash := sha256.Sum256(blockBytes)
	return hex.EncodeToString(hash[:])
}

func validProof(lastProof, proof int) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	guessHash := sha256.Sum256([]byte(guess))
	return hex.EncodeToString(guessHash[:])[:4] == "0000"
}

var blockchain = NewBlockchain()

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/blocks", func(c *gin.Context) {
		response := map[string]interface{}{
			"blocks": blockchain.GetBlocks(),
			"length": len(blockchain.GetBlocks()),
		}
		c.JSON(http.StatusOK, response)
	})

	r.POST("/mine", func(c *gin.Context) {
		var data struct {
			VoterID     string `json:"voter_id"`
			CandidateID string `json:"candidate_id"`
		}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "voter_id e candidate_id is needed"})
			return
		}

		lastBlock := blockchain.LastBlock()
		lastProof := lastBlock.Proof
		proof := blockchain.ProofOfWork(lastProof)

		blockchain.NewTransaction(data.VoterID, data.CandidateID)

		previousHash := Hash(lastBlock)
		block := blockchain.NewBlock(proof, previousHash)

		response := map[string]interface{}{
			"message":       "New Block Created",
			"index":         block.Index,
			"transactions":  block.Transactions,
			"proof":         block.Proof,
			"previous_hash": block.PreviousHash,
		}
		c.JSON(http.StatusOK, response)
	})

	r.GET("/chain", func(c *gin.Context) {
		response := map[string]interface{}{
			"chain":  blockchain.GetBlocks(),
			"length": len(blockchain.GetBlocks()),
		}
		c.JSON(http.StatusOK, response)
	})

	r.Run(":5000")
}
