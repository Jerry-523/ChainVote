package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Estruturas e funções para Block e Transaction
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

func hash(block Block) string {
	blockBytes, _ := json.Marshal(block)
	hash := sha256.Sum256(blockBytes)
	return hex.EncodeToString(hash[:])
}

func proofOfWork(lastProof int) int {
	proof := 0
	for !validProof(lastProof, proof) {
		proof++
	}
	return proof
}

func validProof(lastProof, proof int) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	guessHash := sha256.Sum256([]byte(guess))
	return hex.EncodeToString(guessHash[:])[:4] == "0000"
}

func mine() {
	for {
		resp, err := http.Get("http://localhost:5000/blocks")
		if err != nil {
			fmt.Println("Error getting blocks:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Println("Error decoding response:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		blocks := result["blocks"].([]interface{})
		lastBlock := blocks[len(blocks)-1].(map[string]interface{})
		lastProof := int(lastBlock["proof"].(float64))
		lastHash := lastBlock["previous_hash"].(string)

		proof := proofOfWork(lastProof)
		newBlock := Block{
			Index:        len(blocks) + 1,
			Timestamp:    time.Now().Unix(),
			Transactions: []Transaction{}, // As transações devem ser enviadas pelo servidor principal
			Proof:        proof,
			PreviousHash: lastHash,
		}

		blockBytes, _ := json.Marshal(newBlock)
		resp, err = http.Post("http://localhost:5000/mine", "application/json", bytes.NewBuffer(blockBytes))
		if err != nil {
			fmt.Println("Error mining block:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		

		fmt.Println("Mined a new block:", result)
		time.Sleep(10 * time.Second) // Ajuste o intervalo conforme necessário
	}
}

func main() {
	mine()
}
