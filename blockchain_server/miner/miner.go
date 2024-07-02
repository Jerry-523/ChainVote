// miner/miner.go

package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strconv"
    "time"
    "github.com/Jerry-523/ChainVote/tree/main---blockchian-in-Go/blockchain_server/blockchain"
)

const BlockchainServerURL = "http://localhost:5000"

func getLastBlock() (blockchain.Block, error) {
    resp, err := http.Get(fmt.Sprintf("%s/blocks", BlockchainServerURL))
    if err != nil {
        return blockchain.Block{}, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return blockchain.Block{}, err
    }

    var result struct {
        Blocks []blockchain.Block `json:"blocks"`
    }
    if err := json.Unmarshal(body, &result); err != nil {
        return blockchain.Block{}, err
    }

    if len(result.Blocks) == 0 {
        return blockchain.Block{}, fmt.Errorf("no blocks found")
    }

    return result.Blocks[len(result.Blocks)-1], nil
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

func mineBlock(voterID, candidateID string) {
    lastBlock, err := getLastBlock()
    if err != nil {
        log.Fatalf("Error getting the last block: %v", err)
    }

    lastProof := lastBlock.Proof
    proof := proofOfWork(lastProof)

    data := map[string]string{
        "voter_id":     voterID,
        "candidate_id": candidateID,
    }

    resp, err := http.Post(fmt.Sprintf("%s/mine", BlockchainServerURL), "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Fatalf("Error sending the new block: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Error reading the response body: %v", err)
    }

    fmt.Println("Block mined successfully:", string(body))
}

func main() {
    voterID := "voter123"
    candidateID := "candidate456"

    for {
        mineBlock(voterID, candidateID)
        time.Sleep(10 * time.Second)  // 10 seconds interval between mining blocks
    }
}
