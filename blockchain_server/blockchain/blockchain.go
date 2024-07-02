// blockchain/blockchain.go

package blockchain

import (
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "time"
)

type Transaction struct {
    VoterID     string `json:"voter_id"`
    CandidateID string `json:"candidate_id"`
}

type Block struct {
    Index         int           `json:"index"`
    Timestamp     int64         `json:"timestamp"`
    Transactions  []Transaction `json:"transactions"`
    Proof         int           `json:"proof"`
    PreviousHash  string        `json:"previous_hash"`
}

type Blockchain struct {
    Chain        []Block
    CurrentTransactions []Transaction
    Validators   []string
    MinValidators int
    Voters       map[string]bool
    Candidates   map[string]bool
}

func NewBlockchain() *Blockchain {
    bc := &Blockchain{
        Chain: []Block{},
        CurrentTransactions: []Transaction{},
        Validators: []string{},
        MinValidators: 2,
        Voters: make(map[string]bool),
        Candidates: make(map[string]bool),
    }
    bc.createGenesisBlock()
    return bc
}

func (bc *Blockchain) createGenesisBlock() {
    bc.newBlock(100, "1")  // Genesis block with proof 100 and previous hash "1"
}

func (bc *Blockchain) newBlock(proof int, previousHash string) Block {
    block := Block{
        Index:         len(bc.Chain) + 1,
        Timestamp:     time.Now().Unix(),
        Transactions:  bc.CurrentTransactions,
        Proof:         proof,
        PreviousHash:  previousHash,
    }
    bc.CurrentTransactions = []Transaction{}
    bc.Chain = append(bc.Chain, block)
    return block
}

func (bc *Blockchain) newTransaction(voterID, candidateID string) error {
    if _, exists := bc.Voters[voterID]; exists {
        return errors.New("voter has already voted")
    }
    if _, exists := bc.Candidates[candidateID]; !exists {
        return errors.New("candidate does not exist")
    }
    bc.CurrentTransactions = append(bc.CurrentTransactions, Transaction{VoterID: voterID, CandidateID: candidateID})
    bc.Voters[voterID] = true
    return nil
}

func (bc *Blockchain) lastBlock() Block {
    return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) hash(block Block) string {
    blockString, _ := json.Marshal(block)
    hash := sha256.Sum256(blockString)
    return hex.EncodeToString(hash[:])
}

func (bc *Blockchain) proofOfWork(lastProof int) int {
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

func (bc *Blockchain) addValidator(validatorID string) error {
    if len(bc.Validators) >= bc.MinValidators {
        return errors.New("maximum number of validators reached")
    }
    bc.Validators = append(bc.Validators, validatorID)
    return nil
}

func (bc *Blockchain) registerVoter(voterID string) {
    bc.Voters[voterID] = true
}

func (bc *Blockchain) registerCandidate(candidateID string) {
    bc.Candidates[candidateID] = true
}

func (bc *Blockchain) getBlocks() []Block {
    return bc.Chain
}

func (bc *Blockchain) getChain() []Block {
    return bc.Chain
}

func (bc *Blockchain) validateBlock(block Block) bool {
    if !bc.validProof(bc.lastBlock().Proof, block.Proof) {
        return false
    }
    for _, tx := range block.Transactions {
        if _, exists := bc.Voters[tx.VoterID]; !exists {
            return false
        }
        if _, exists := bc.Candidates[tx.CandidateID]; !exists {
            return false
        }
    }
    return true
}

func (bc *Blockchain) getUnvalidatedBlocks() []Block {
    // This function will return blocks that need to be validated by validators
    return []Block{}  // Placeholder for demonstration
}

