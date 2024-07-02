// validator/validator.go

package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "github.com/Jerry-523/ChainVote/tree/main---blockchian-in-Go/blockchain_server/blockchain"
)

const BlockchainServerURL = "http://localhost:5000"

func getUnvalidatedBlocks() ([]blockchain.Block, error) {
    resp, err := http.Get(fmt.Sprintf("%s/unvalidated-blocks", BlockchainServerURL))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var result struct {
        Blocks []blockchain.Block `json:"blocks"`
    }
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result.Blocks, nil
}

func validateBlock(block blockchain.Block) {
    jsonData, err := json.Marshal(block)
    if err != nil {
        log.Fatalf("Error marshalling block: %v", err)
    }

    resp, err := http.Post(fmt.Sprintf("%s/validate-block", BlockchainServerURL), "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Fatalf("Error validating the block: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Error reading the response body: %v", err)
    }

    fmt.Println("Block validated successfully:", string(body))
}

func main() {
    for {
        blocks, err := getUnvalidatedBlocks()
        if err != nil {
            log.Fatalf("Error getting unvalidated blocks: %v", err)
        }

        for _, block := range blocks {
            validateBlock(block)
        }
    }
}
