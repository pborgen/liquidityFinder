package walletFinderService

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
	"github.com/pborgen/liquidityFinder/internal/testhelper"
	"github.com/pborgen/liquidityFinder/myConst"
	"github.com/rs/zerolog/log"
)

func init() {
	testhelper.Setup()
}

func TestBla(t *testing.T) {
	// hexChars := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	// client := blockchainclient.GetHttpClient()

	privateKeyInt := big.NewInt(1000)
    one := big.NewInt(1)
    maxIterations := int64(1000000) // Limit iterations for demonstration (e.g., first 1000 keys)


	for i := int64(0); i < maxIterations; i++ {
        // Convert to 32-byte slice (pad left with zeros if needed)
        privateKeyBytes := privateKeyInt.Bytes()
        if len(privateKeyBytes) < 32 {
            padded := make([]byte, 32)
            copy(padded[32-len(privateKeyBytes):], privateKeyBytes)
            privateKeyBytes = padded
        }

        // Parse as ECDSA private key
        privateKey, err := crypto.ToECDSA(privateKeyBytes)
        if err != nil {
            log.Printf("Invalid private key at iteration %d: %v", i, err)
			privateKeyInt.Add(privateKeyInt, one)
            continue
        }

        // Derive public key and address
        publicKey := privateKey.Public().(*ecdsa.PublicKey)
        address := crypto.PubkeyToAddress(*publicKey)

        checkBalance(address)

        // Convert private key to hex string
        privateKeyHex := fmt.Sprintf("%064x", privateKey.D) // D is the private key scalar

        // Print details
        fmt.Printf("Private Key: 0x%s\n", privateKeyHex)
        fmt.Printf("Address: %s\n", address.Hex())
       
        fmt.Println("---")

        // Increment private key
        privateKeyInt.Add(privateKeyInt, one)
    }

	// for _, char := range hexChars {
	// 	privateKeyHex := fmt.Sprintf("%s%s", string(char), string(char))
	// 	log.Info().Msgf("Private Key Hex: %s", privateKeyHex)


	// 	balance, err := client.BalanceAt(context.Background(), common.HexToAddress("0x0000000000000000000000000000000000000000"), nil)
	// 	if err != nil {
	// 		log.Fatal().Msgf("Failed to get balance: %v", err)
	// 	}
	// 	log.Info().Msgf("Balance: %s", balance.String())

	// 	if balance.Cmp(big.NewInt(0)) > 0 {
	// 		log.Info().Msgf("Balance is greater than 0: %s", balance.String())
	// 	}
	// }

	privateKeyHex := "0000000000000000000000000000000000000000000000000000000000000001"

    // Parse the private key
    privateKey, err := crypto.HexToECDSA(privateKeyHex)
    if err != nil {
        log.Fatal().Msgf("Failed to parse private key: %v", err)
    }

    // Derive the public key
    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal().Msgf("Error casting public key to ECDSA")
    }

    // Convert to bytes (uncompressed form, 65 bytes: 04 || X || Y)
    publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
    log.Info().Msgf("Public Key (hex): %x\n", publicKeyBytes)

    // Optional: Derive Ethereum address from public key
    address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
    fmt.Printf("Ethereum Address: %s\n", address)

	log.Info().Msgf("one ddddpls: %s", myConst.GetOneWplsBigint().String())
}

func checkBalance(address common.Address) {
	client := blockchainclient.GetHttpClient()
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal().Msgf("Failed to get balance: %v", err)
	}

	minBalance := big.NewInt(0).Mul(big.NewInt(100), big.NewInt(1000000000000000000))

	if balance.Cmp(minBalance) > 0 {
		log.Info().Msgf("Balance is greater than 0: %s", balance.String())
		log.Info().Msgf("Address: %s", address.Hex())
	}
}