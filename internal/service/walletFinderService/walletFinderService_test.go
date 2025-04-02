package walletFinderService

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
	"github.com/pborgen/liquidityFinder/internal/testhelper"
	"github.com/rs/zerolog/log"
)

func init() {
	testhelper.Setup()
}

func TestBla(t *testing.T) {
	numWorkers := 10
	maxIterations := int64(1000000)
	batchSize := int64(1000)

	// Channel for workers to signal completion
	var wg sync.WaitGroup
	errChan := make(chan error, numWorkers)

	// Create a channel to distribute work
	workChan := make(chan int64, numWorkers)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Each worker gets its own private key to increment
			privateKeyInt := big.NewInt(0)
			one := big.NewInt(1)

			for startIndex := range workChan {
				endIndex := startIndex + batchSize
				if endIndex > maxIterations {
					endIndex = maxIterations
				}

				// Set initial private key for this batch
				privateKeyInt.SetInt64(startIndex)

				for i := startIndex; i < endIndex; i++ {
					// Convert to 32-byte slice
					privateKeyBytes := privateKeyInt.Bytes()
					if len(privateKeyBytes) < 32 {
						padded := make([]byte, 32)
						copy(padded[32-len(privateKeyBytes):], privateKeyBytes)
						privateKeyBytes = padded
					}

					// Parse as ECDSA private key
					privateKey, err := crypto.ToECDSA(privateKeyBytes)
					if err != nil {
						log.Debug().
							Int("worker", workerID).
							Int64("iteration", i).
							Err(err).
							Msg("Invalid private key")
						privateKeyInt.Add(privateKeyInt, one)
						continue
					}

					// Derive public key and address
					publicKey := privateKey.Public().(*ecdsa.PublicKey)
					address := crypto.PubkeyToAddress(*publicKey)

					// Check balance
					checkBalance(address)

					// Convert private key to hex string for logging
					if i%1000 == 0 { // Log less frequently
						privateKeyHex := fmt.Sprintf("%064x", privateKey.D)
						log.Debug().
							Int("worker", workerID).
							Int64("iteration", i).
							Str("privateKey", privateKeyHex).
							Str("address", address.Hex()).
							Msg("Processing")
					}

					// Increment private key
					privateKeyInt.Add(privateKeyInt, one)

				}
			}
		}(i)
	}

	// Distribute work to workers
	go func() {
		for i := int64(0); i < maxIterations; i += batchSize {
			workChan <- i
		}
		close(workChan)
	}()

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for any errors
	for err := range errChan {
		if err != nil {
			t.Errorf("Worker error: %v", err)
		}
	}
}

func checkBalance(address common.Address) {
	client := blockchainclient.GetPublicHttpClient()
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Error().
			Err(err).
			Str("address", address.Hex()).
			Msg("Failed to get balance")
		return // Don't fatal in a goroutine
	}

	minBalance := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1000000000000000000))

	if balance.Cmp(minBalance) > 0 {
		log.Info().
			Str("balance", balance.String()).
			Str("address", address.Hex()).
			Msg("Found address with large balance!")
	}
}