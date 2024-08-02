// server.go

package main

import (
    "net/http"
    "blockchain_server/blockchain" // Importando o pacote blockchain
    "github.com/gin-gonic/gin"
)

var (
    votes             = make(map[int]map[string]bool)
    validators        = make(map[string]bool)
    minimumValidators = 2
)

func getBlocks(bc *blockchain.Blockchain) gin.HandlerFunc {
    return func(c *gin.Context) {
        blocks := bc.GetBlocks()
        c.JSON(http.StatusOK, gin.H{
            "blocks": blocks,
            "length": len(blocks),
        })
    }
}

func mineHandler(bc *blockchain.Blockchain) gin.HandlerFunc {
    return func(c *gin.Context) {
        var data struct {
            VoterID     string `json:"voter_id"`
            CandidateID string `json:"candidate_id"`
        }

        if err := c.BindJSON(&data); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de transação inválidos"})
            return
        }

        if data.VoterID == "" || data.CandidateID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Campos voter_id e candidate_id são obrigatórios"})
            return
        }

        bc.NewTransaction(data.VoterID, data.CandidateID)

        lastBlock := bc.LastBlock()
        lastProof := lastBlock.Proof
        proof := bc.ProofOfWork(lastProof)

        previousHash := bc.Hash(lastBlock)
        block := bc.NewBlock(proof, previousHash)

        c.JSON(http.StatusOK, gin.H{
            "message":        "Novo bloco criado",
            "index":          block.Index,
            "transactions":   block.Transactions,
            "proof":          block.Proof,
            "previous_hash":  block.PreviousHash,
        })
    }
}

func fullChainHandler(bc *blockchain.Blockchain) gin.HandlerFunc {
    return func(c *gin.Context) {
        chain := bc.GetChain()
        c.JSON(http.StatusOK, gin.H{
            "chain":  chain,
            "length": len(chain),
        })
    }
}

func addValidatorHandler(c *gin.Context) {
    validatorID := c.Query("validator_id")
    if validatorID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "validator_id é obrigatório"})
        return
    }

    validators[validatorID] = true
    c.JSON(http.StatusOK, gin.H{"message": "Validador adicionado com sucesso"})
}

func voteHandler(bc *blockchain.Blockchain) gin.HandlerFunc {
    return func(c *gin.Context) {
        validatorID := c.GetHeader("Validator-ID")
        if validatorID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Header Validator-ID é obrigatório"})
            return
        }

        if _, ok := validators[validatorID]; !ok {
            c.JSON(http.StatusForbidden, gin.H{"error": "Validador não autorizado"})
            return
        }

        var voteData struct {
            BlockIndex int  `json:"block_index"`
            Vote       bool `json:"vote"`
        }

        if err := c.BindJSON(&voteData); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de voto inválidos"})
            return
        }

        block, exists := bc.GetBlock(voteData.BlockIndex)
        if !exists {
            c.JSON(http.StatusNotFound, gin.H{"error": "Bloco não encontrado"})
            return
        }

        if votes[voteData.BlockIndex] == nil {
            votes[voteData.BlockIndex] = make(map[string]bool)
        }

        votes[voteData.BlockIndex][validatorID] = voteData.Vote

        if isBlockConsensusAchieved(voteData.BlockIndex) {
            if isBlockAccepted(voteData.BlockIndex) {
                if bc.AddBlock(block) {
                    c.JSON(http.StatusOK, gin.H{"message": "Bloco adicionado com sucesso"})
                } else {
                    c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao adicionar o bloco"})
                }
            } else {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Bloco rejeitado pela maioria dos validadores"})
            }
        } else {
            c.JSON(http.StatusAccepted, gin.H{"message": "Voto registrado, aguardando consenso"})
        }
    }
}

func isBlockConsensusAchieved(blockIndex int) bool {
    return len(votes[blockIndex]) >= minimumValidators
}

func isBlockAccepted(blockIndex int) bool {
    positiveVotes := 0
    for _, vote := range votes[blockIndex] {
        if vote {
            positiveVotes++
        }
    }
    return positiveVotes > len(votes[blockIndex])/2
}
