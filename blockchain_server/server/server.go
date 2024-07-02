// server/server.go

package main

import (
    "encoding/json"
    "net/http"
    "github.com/Jerry-523/ChainVote/tree/main---blockchian-in-Go/blockchain_server/blockchain"
    "github.com/gorilla/mux"
)

var bc = blockchain.NewBlockchain()

func getBlocksHandler(w http.ResponseWriter, r *http.Request) {
    blocks := bc.GetBlocks()
    response := map[string]interface{}{
        "blocks": blocks,
    }
    json.NewEncoder(w).Encode(response)
}

func mineHandler(w http.ResponseWriter, r *http.Request) {
    var data map[string]string
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    voterID := data["voter_id"]
    candidateID := data["candidate_id"]

    if err := bc.NewTransaction(voterID, candidateID); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    lastBlock := bc.LastBlock()
    lastProof := lastBlock.Proof
    proof := bc.ProofOfWork(lastProof)
    previousHash := bc.Hash(lastBlock)
    block := bc.NewBlock(proof, previousHash)

    response := map[string]interface{}{
        "message":        "New block created",
        "index":          block.Index,
        "transactions":   block.Transactions,
        "proof":          block.Proof,
        "previous_hash":  block.PreviousHash,
    }
    json.NewEncoder(w).Encode(response)
}

func fullChainHandler(w http.ResponseWriter, r *http.Request) {
    chain := bc.GetChain()
    response := map[string]interface{}{
        "chain": chain,
        "length": len(chain),
    }
    json.NewEncoder(w).Encode(response)
}

func addValidatorHandler(w http.ResponseWriter, r *http.Request) {
    var data map[string]string
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    validatorID := data["validator_id"]

    if err := bc.AddValidator(validatorID); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    response := map[string]interface{}{
        "message":    "Validator added successfully",
        "validators": bc.Validators,
    }
    json.NewEncoder(w).Encode(response)
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
    var data map[string]string
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    voterID := data["voter_id"]
    candidateID := data["candidate_id"]

    if err := bc.NewTransaction(voterID, candidateID); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    response := map[string]string{
        "message": "Vote recorded successfully",
    }
    json.NewEncoder(w).Encode(response)
}

func validateBlockHandler(w http.ResponseWriter, r *http.Request) {
    var block blockchain.Block
    if err := json.NewDecoder(r.Body).Decode(&block); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if !bc.ValidateBlock(block) {
        http.Error(w, "Block is invalid", http.StatusBadRequest)
        return
    }

    response := map[string]string{
        "message": "Block validated successfully",
    }
    json.NewEncoder(w).Encode(response)
}

func getUnvalidatedBlocksHandler(w http.ResponseWriter, r *http.Request) {
    blocks := bc.GetUnvalidatedBlocks()
    response := map[string]interface{}{
        "blocks": blocks,
    }
    json.NewEncoder(w).Encode(response)
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/blocks", getBlocksHandler).Methods("GET")
    r.HandleFunc("/mine", mineHandler).Methods("POST")
    r.HandleFunc("/chain", fullChainHandler).Methods("GET")
    r.HandleFunc("/add-validator", addValidatorHandler).Methods("POST")
    r.HandleFunc("/vote", voteHandler).Methods("POST")
    r.HandleFunc("/validate-block", validateBlockHandler).Methods("POST")
    r.HandleFunc("/unvalidated-blocks", getUnvalidatedBlocksHandler).Methods("GET")

    http.ListenAndServe(":5000", r)
}
