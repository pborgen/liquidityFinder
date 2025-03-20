package blockchainutil

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
	"github.com/pborgen/liquidityFinder/internal/myUtil"
	"github.com/rs/zerolog/log"
)

// RPCResponse represents the full JSON-RPC response
type RPCResponse struct {
    Jsonrpc string          `json:"jsonrpc"`
    ID      int             `json:"id"`
    Result  TraceCallResult `json:"result"`
}

// TraceCallResult represents the debug_traceCall result
type TraceCallResult struct {
    Gas         int64         `json:"gas"`         // Changed to int64 since it's a number
    Failed      bool          `json:"failed"`
    ReturnValue string        `json:"returnValue"`
    StructLogs  []StructLog   `json:"structLogs"`
}

// StructLog represents each step in the execution trace
type StructLog struct {
    Pc       uint64                 `json:"pc"`
    Op       string                 `json:"op"`
    Gas      uint64                 `json:"gas"`
    GasCost  uint64                 `json:"gasCost"`
    Depth    int                    `json:"depth"`
    Stack    []string               `json:"stack"`
    Storage  map[string]string      `json:"storage,omitempty"` // Added for storage changes
}

type TransactionError struct {
	Err                     error
	RevertReason            string
	DidExecuteTransaction   bool
	IsTransactionSuccessful bool
}

func (w *TransactionError) Error() string {
	return fmt.Sprintf("DidExecuteTransaction: %t, IsTransactionSuccessful: %t,  %v", w.DidExecuteTransaction, w.IsTransactionSuccessful, w.Err)
}

func GetCurrentBlockNumber() uint64 {
	header, err := GetCurrentBlockHeader()
	var blockNumber uint64 = 0

	if err != nil {
		log.Err(err)
	} else {

		blockNumber = uint64(header.Number.Uint64())
	}

	return blockNumber
}

func GetCurrentBlockNumberAsUInt64() uint64 {
	return GetCurrentBlockNumber()
}

func GetCurrentBlockNumberAsBigInt() *big.Int {
	header, err := GetCurrentBlockHeader()
	var blockNumber big.Int

	if err != nil {
		log.Err(err)
	} else {

		blockNumber = *header.Number
	}

	return &blockNumber
}

func GetCurrentBlockTimestamp() (uint64, error) {
	header, err := GetCurrentBlockHeader()
	if err != nil {
		return 0, err
	}
	return header.Time, nil
}

func GetCurrentBlockHeader() (*types.Header, error) {
	client := blockchainclient.GetHttpClient()
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func IsContractAddress(address common.Address) (bool, error) {
	var client = blockchainclient.GetHttpClient()
	// Get the code at the address
	bytecode, err := client.CodeAt(context.Background(), address, nil) // nil is the latest block
	if err != nil {
		return false, err
	}

	// If the bytecode length is greater than zero, it's a contract
	return len(bytecode) > 0, nil
}

func GetTransactionResults(opts *bind.TransactOpts, transaction *types.Transaction) (bool, error) {
	client := blockchainclient.GetHttpClient()

	receipt, err := bind.WaitMined(context.Background(), client, transaction)

	if err != nil {
		log.Log().Msgf("Error waiting for transaction to be mined: %v", err)

		myError := &TransactionError{
			DidExecuteTransaction: true,
			Err:                   err,
		}

		return false, myError
	}

	isSuccessful, err := myUtil.IsTransactionSuccess(receipt)

	if err != nil {
		myError := &TransactionError{
			DidExecuteTransaction:   true,
			IsTransactionSuccessful: isSuccessful,
			Err:                     err,
		}

		return false, myError
	}

	if isSuccessful {
		log.Info().Msgf("Transaction successful: %s", receipt.TxHash.String())
	} else {
		revertReason, err := blockchainclient.GetRevertReason(opts.From, transaction)

		myError := &TransactionError{
			DidExecuteTransaction:   true,
			RevertReason:            revertReason,
			IsTransactionSuccessful: isSuccessful,
			Err:                     err,
		}

		return false, myError
	}
	

	return true, nil
}

func GetTransactionSuggestedFees() (gasFeeCap *big.Int, gasTipCap *big.Int, baseFee *big.Int, err error) {
	client := blockchainclient.GetHttpClient()

	// Fetch the latest block header
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Error().Msgf("Failed to get latest block header: %v", err)
		return nil, nil, nil, err
	}

	// Get the current base fee
	baseFee = header.BaseFee

	// Suggest a gas tip cap
	gasTipCap, err = client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Error().Msgf("Failed to suggest gas tip cap: %v", err)
		return nil, nil, nil, err
	}

	buffer := big.NewInt(2)
	gasFeeCap = new(big.Int).Add(baseFee, gasTipCap)
	gasFeeCap.Mul(gasFeeCap, buffer)


	return gasFeeCap, gasTipCap, baseFee, nil
}

func GetNextNonce(address common.Address) (*big.Int, error) {
	client := blockchainclient.GetHttpClient()

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}

	return big.NewInt(int64(nonce)), nil
}


func SetupAuthForTransaction(auth *bind.TransactOpts, addExtraGasTip bool) error {
	gasFeeCap, suggestedGasTipCap, _, err := GetTransactionSuggestedFees()
	if err != nil {
		log.Fatal().Msgf("Error getting transaction suggested fees: %v", err)
	}

	nonce, err := GetNextNonce(auth.From)
	log.Info().Msgf("Nonce: %s", nonce.String())
	if err != nil {
		log.Fatal().Msgf("Error getting next nonce: %v", err)
		return err
	}

	if addExtraGasTip {
		// add 1% extra gas tip
		onePercentOfSuggestedGasTipCap := new(big.Int).Div(suggestedGasTipCap, big.NewInt(100))
		suggestedGasTipCap = new(big.Int).Add(suggestedGasTipCap, onePercentOfSuggestedGasTipCap)
	}
	
	auth.Nonce = nonce
	auth.GasPrice = nil
	auth.GasFeeCap = gasFeeCap
	auth.GasTipCap = suggestedGasTipCap
	auth.GasLimit = uint64(0) // Auto estimate the gas limit
	
	
	return nil
}

func GetCostInPlsForTransaction(gasLimit *big.Int, gasPrice *big.Int, gasTipCap *big.Int) (*big.Int) {

	gasPlusTip := new(big.Int).Add(gasPrice, gasTipCap)

	costInPls := new(big.Int).Mul(gasLimit, gasPlusTip)

	return costInPls
}

func GetAuthAccount(privateKeyString string, chainId int) (*bind.TransactOpts, error) {

	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return nil, err
	}


	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(chainId)))
	
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func GeneratePublicAndPrivateKey() (publicKeyString string, privateKeyString string) {

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	publicKey := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	// Convert the private key to a hexadecimal string
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyString = fmt.Sprintf("%x", privateKeyBytes)

	return publicKey, privateKeyString
}



func DebugTraceCall(
	from common.Address,
	to common.Address, 
	data []byte,
) (string, error) {
	url := "https://docs-demo.quiknode.pro/"

	hexString := "0x" + hex.EncodeToString(data)

	log.Info().Msgf("Hex string: %s", hexString)
    // Define the payload structure
    payload := struct {
        Method  string        `json:"method"`
        Params  []interface{} `json:"params"`
        ID      int           `json:"id"`
        Jsonrpc string        `json:"jsonrpc"`
    }{
        Method: "debug_traceCall",
        Params: []interface{}{
            map[string]interface{}{
                "from": from,
                "to":   to,
                "data": hexString,
            },
            "latest",
        },
        ID:      1,
        Jsonrpc: "2.0",
    }

    // Convert payload to JSON
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }

    // Create HTTP request
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
    if err != nil {
        return "", err
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")

    // Create HTTP client and execute request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
    
        return "", err
    }
    defer resp.Body.Close()

    // Read and print response
    buf := new(bytes.Buffer)
    _, err = buf.ReadFrom(resp.Body)
    if err != nil {
        fmt.Printf("Error reading response: %v\n", err)
        return "", err
    }
	
	var response RPCResponse
    err = json.Unmarshal([]byte(buf.String()), &response)
    if err != nil {
        fmt.Printf("Error unmarshaling JSON: %v\n", err)
        return "", err
    }

    // Convert back to formatted JSON
    jsonOutput, err := json.Marshal(response)
    if err != nil {
        return "", err
    }

	

	return string(jsonOutput), nil
}

func TraceCall(
	from common.Address,
	to common.Address, 
	data []byte,
) (string, error) {
	url := "https://docs-demo.quiknode.pro/"

	hexString := "0x" + hex.EncodeToString(data)
	
	log.Info().Msgf("Hex string: %s", hexString)
    // Define the payload structure
    payload := struct {
        Method  string        `json:"method"`
        Params  []interface{} `json:"params"`
        ID      int           `json:"id"`
        Jsonrpc string        `json:"jsonrpc"`
    }{
        Method: "trace_call",
        Params: []interface{}{
            map[string]interface{}{
                "from": from,
                "to":   to,
                "data": hexString,
            },
			[]string{"trace"},
            "latest",
        },
        ID:      1,
        Jsonrpc: "2.0",
    }

    // Convert payload to JSON
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }

    // Create HTTP request
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
    if err != nil {
        return "", err
    }

    // Set headers
    req.Header.Set("Content-Type", "application/json")

    // Create HTTP client and execute request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
    
        return "", err
    }
    defer resp.Body.Close()

    // Read and print response
    buf := new(bytes.Buffer)
    _, err = buf.ReadFrom(resp.Body)
    if err != nil {
        fmt.Printf("Error reading response: %v\n", err)
        return "", err
    }
	
	var response RPCResponse
    err = json.Unmarshal([]byte(buf.String()), &response)
    if err != nil {
        fmt.Printf("Error unmarshaling JSON: %v\n", err)
        return "", err
    }

    // Convert back to formatted JSON
    jsonOutput, err := json.Marshal(response)
    if err != nil {
        return "", err
    }

	

	return string(jsonOutput), nil
}
