package transferEventGather

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	blockchainutil "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainutil"
)

// TransferEvent represents a token transfer event
type TransferEvent struct {
    From   common.Address
    To     common.Address
    Value  *big.Int
    TxHash common.Hash
    Block  uint64
}

// TokenEventTracker handles token event tracking
type TokenEventTracker struct {
    client *ethclient.Client
    abi    abi.ABI
}

// Standard ERC20 Transfer event signature
const ERC20ABI = `[{
    "anonymous": false,
    "inputs": [
        {
            "indexed": true,
            "name": "from",
            "type": "address"
        },
        {
            "indexed": true,
            "name": "to",
            "type": "address"
        },
        {
            "indexed": false,
            "name": "value",
            "type": "uint256"
        }
    ],
    "name": "Transfer",
    "type": "event"
}]`

func NewTokenEventTracker(client *ethclient.Client) (*TokenEventTracker, error) {


    contractAbi, err := abi.JSON(strings.NewReader(ERC20ABI))
    if err != nil {
        return nil, fmt.Errorf("failed to parse ABI: %v", err)
    }

    return &TokenEventTracker{
        client: client,
        abi:    contractAbi,
    }, nil
}



func (t *TokenEventTracker) GetAllTransferEventsForBlockRange(
    ctx context.Context,
    fromBlock uint64,
    toBlock uint64,
) ([]TransferEvent, error) {

    // Create the query
    query := ethereum.FilterQuery{
        FromBlock: new(big.Int).SetUint64(fromBlock),
        ToBlock:   new(big.Int).SetUint64(toBlock),
        Topics: [][]common.Hash{
            {t.abi.Events["Transfer"].ID}, // Transfer event signature
        },
    }

    // Get logs
    logs, err := t.client.FilterLogs(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to filter logs: %v", err)
    }

    // Process logs
    var events []TransferEvent
    for _, log := range logs {
        event, err := t.parseTransferEvent(log)
        if err != nil {
            continue // Skip failed parsing
        }
        events = append(events, event)
    }

    return events, nil
}

// GetTransferEvents retrieves transfer events within a block range
func (t *TokenEventTracker) GetAllTransferEventsForLatestBlock(
    ctx context.Context,
) ([]TransferEvent, uint64, error) {

    currentBlock := blockchainutil.GetCurrentBlockNumber()
    events, err := t.GetAllTransferEventsForBlockRange(ctx, currentBlock, currentBlock)
    if err != nil {
        return nil, 0, err
    }
    
    return events, currentBlock, nil
}

// GetTransferEvents retrieves transfer events within a block range
func (t *TokenEventTracker) GetTransferEvents(
    ctx context.Context,
    tokenAddress common.Address,
    fromBlock, toBlock *big.Int,
) ([]TransferEvent, error) {
    // Create the query
    // Define the Transfer event signature
	transferEventSignature := []byte("Transfer(address,address,uint256)")
	transferEventHash := common.BytesToHash(crypto.Keccak256(transferEventSignature))

    
    query := ethereum.FilterQuery{
        FromBlock: fromBlock,
        ToBlock:   toBlock,
        Addresses: []common.Address{tokenAddress},
        Topics: [][]common.Hash{
            {transferEventHash},
            //{t.abi.Events["Transfer"].ID}, // Transfer event signature
        },
    }

    // Get logs
    logs, err := t.client.FilterLogs(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to filter logs: %v", err)
    }

    // Process logs
    var events []TransferEvent
    for _, log := range logs {
        event, err := t.parseTransferEvent(log)
        if err != nil {
            continue // Skip failed parsing
        }
        events = append(events, event)
    }

    return events, nil
}

// GetLatestTransfers gets the most recent transfers
func (t *TokenEventTracker) GetLatestTransfers(
    ctx context.Context,
    tokenAddress common.Address,
    numBlocks int64,
) ([]TransferEvent, error) {
    // Get current block
    header, err := t.client.HeaderByNumber(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get latest block: %v", err)
    }

    currentBlock := header.Number
    fromBlock := new(big.Int).Sub(currentBlock, big.NewInt(numBlocks))
    if fromBlock.Cmp(big.NewInt(0)) < 0 {
        fromBlock = big.NewInt(0)
    }

    return t.GetTransferEvents(ctx, tokenAddress, fromBlock, currentBlock)
}

// GetTransfersByAddress gets transfers for a specific address
func (t *TokenEventTracker) GetTransfersByAddress(
    ctx context.Context,
    tokenAddress common.Address,
    userAddress common.Address,
    fromBlock, toBlock *big.Int,
) ([]TransferEvent, error) {
    // Create topics for filtering
    addressTopic := common.BytesToHash(userAddress.Bytes())
    query := ethereum.FilterQuery{
        FromBlock: fromBlock,
        ToBlock:   toBlock,
        Addresses: []common.Address{tokenAddress},
        Topics: [][]common.Hash{
            {t.abi.Events["Transfer"].ID},
            {addressTopic}, // From address
            {addressTopic}, // To address
        },
    }

    logs, err := t.client.FilterLogs(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to filter logs: %v", err)
    }

    var events []TransferEvent
    for _, log := range logs {
        event, err := t.parseTransferEvent(log)
        if err != nil {
            continue
        }
        events = append(events, event)
    }

    return events, nil
}

// parseTransferEvent parses a log into a TransferEvent
func (t *TokenEventTracker) parseTransferEvent(log types.Log) (TransferEvent, error) {
    event := TransferEvent{
        TxHash: log.TxHash,
        Block:  log.BlockNumber,
    }

    // Parse the From address (first topic after event signature)
    if len(log.Topics) >= 2 {
        event.From = common.HexToAddress(log.Topics[1].Hex())
    }

    // Parse the To address (second topic after event signature)
    if len(log.Topics) >= 3 {
        event.To = common.HexToAddress(log.Topics[2].Hex())
    }

    // Parse the value
    event.Value = new(big.Int).SetBytes(log.Data)

    return event, nil
}

// StreamTransfers streams transfer events in real-time
func (t *TokenEventTracker) StreamTransfers(
    ctx context.Context,
    tokenAddress common.Address,
    eventChan chan<- TransferEvent,
) error {
    query := ethereum.FilterQuery{
        Addresses: []common.Address{tokenAddress},
        Topics: [][]common.Hash{
            {t.abi.Events["Transfer"].ID},
        },
    }

    logs := make(chan types.Log)
    sub, err := t.client.SubscribeFilterLogs(ctx, query, logs)
    if err != nil {
        return fmt.Errorf("failed to subscribe to logs: %v", err)
    }

    go func() {
        defer sub.Unsubscribe()
        for {
            select {
            case <-ctx.Done():
                return
            case err := <-sub.Err():
                fmt.Printf("Subscription error: %v\n", err)
                return
            case log := <-logs:
                event, err := t.parseTransferEvent(log)
                if err != nil {
                    continue
                }
                eventChan <- event
            }
        }
    }()

    return nil
}