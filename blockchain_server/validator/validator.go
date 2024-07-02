package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func voteOnBlock(blockIndex int, vote bool) error {
	validatorID := "validator1" // O ID do validador deve ser configurado ou passado como argumento

	data := map[string]interface{}{
		"block_index": blockIndex,
		"vote":        vote,
	}

	dataBytes, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "http://localhost:5000/vote", bytes.NewBuffer(dataBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Validator-ID", validatorID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to vote, status code: %d", resp.StatusCode)
	}

	return nil
}

func getStatus() (map[string]interface{}, error) {
	resp, err := http.Get("http://localhost:5000/chain")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: validator [command]")
		fmt.Println("Commands:")
		fmt.Println("  vote [block_index] [true|false]  Vote on a block (true for accept, false for reject)")
		fmt.Println("  status                          Get the current status of the blockchain")
		return
	}

	switch os.Args[1] {
	case "vote":
		if len(os.Args) != 4 {
			fmt.Println("Usage: validator vote [block_index] [true|false]")
			return
		}
		blockIndex := os.Args[2]
		vote := os.Args[3]
		voteBool := vote == "true"
		blockIndexInt := 0
		fmt.Sscanf(blockIndex, "%d", &blockIndexInt)
		if err := voteOnBlock(blockIndexInt, voteBool); err != nil {
			fmt.Println("Error voting on block:", err)
		} else {
			fmt.Println("Vote cast successfully.")
		}

	case "status":
		status, err := getStatus()
		if err != nil {
			fmt.Println("Error getting status:", err)
		} else {
			fmt.Println("Blockchain status:", status)
		}

	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
