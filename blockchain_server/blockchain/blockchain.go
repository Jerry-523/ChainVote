// blockchain/blockchain.go

package blockchain

import (
    "crypto/sha256"
    "encoding/hex"
    "time"
	"fmt"
    "encoding/json"
)

type Block struct {
    Index        int
    Timestamp    time.Time
    Transactions []Transaction
    Proof        int
    PreviousHash string
}

type Transaction struct {
    VoterID     string
    CandidateID string
}

type Blockchain struct {
    Chain               []Block
    CurrentTransactions []Transaction
}

func NewBlockchain() *Blockchain {
    bc := &Blockchain{}
    bc.NewBlock(100, "1")
    return bc
}

func (bc *Blockchain) NewBlock(proof int, previousHash string) Block {
    block := Block{
        Index:        len(bc.Chain) + 1,
        Timestamp:    time.Now(),
        Transactions: bc.CurrentTransactions,
        Proof:        proof,
        PreviousHash: previousHash,
    }

    bc.CurrentTransactions = nil
    bc.Chain = append(bc.Chain, block)
    return block
}

func (bc *Blockchain) NewTransaction(voterID, candidateID string) int {
    bc.CurrentTransactions = append(bc.CurrentTransactions, Transaction{VoterID: voterID, CandidateID: candidateID})
    return bc.LastBlock().Index + 1
}

func (bc *Blockchain) LastBlock() Block {
    return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) ProofOfWork(lastProof int) int {
    proof := 0
    for !bc.validProof(lastProof, proof) {
        proof++
    }
    return proof
}

func (bc *Blockchain) validProof(lastProof, proof int) bool {
    guess := fmt.Sprintf("%d%d", lastProof, proof)
    guessHash := sha256.Sum256([]byte(guess))
    return hex.EncodeToString(guessHash[:])[:4] == "0000"
}

func (bc *Blockchain) Hash(block Block) string {
    blockBytes, _ := json.Marshal(block)
    hash := sha256.Sum256(blockBytes)
    return hex.EncodeToString(hash[:])
}

func (bc *Blockchain) GetBlocks() []Block {
    return bc.Chain
}

func (bc *Blockchain) GetChain() []Block {
    return bc.Chain
}

func (bc *Blockchain) GetBlock(index int) (Block, bool) {
    if index < 0 || index >= len(bc.Chain) {
        return Block{}, false
    }
    return bc.Chain[index], true
}

func (bc *Blockchain) AddBlock(block Block) bool {
    if bc.validBlock(block, bc.LastBlock()) {
        bc.Chain = append(bc.Chain, block)
        return true
    }
    return false
}

func (bc *Blockchain) validBlock(block, lastBlock Block) bool {
    if block.PreviousHash != bc.Hash(lastBlock) {
        return false
    }
    if !bc.validProof(lastBlock.Proof, block.Proof) {
        return false
    }
    return true
}
