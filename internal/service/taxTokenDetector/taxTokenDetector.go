package taxTokenDetector

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
	"github.com/pborgen/liquidityFinder/internal/blockchain/blockchainutil"
	"github.com/pborgen/liquidityFinder/internal/database/model/erc20"
	"github.com/pborgen/liquidityFinder/internal/database/model/pair"
	"github.com/pborgen/liquidityFinder/internal/service/transferEventGather"
	"github.com/rs/zerolog/log"
)

// Common ERC20 function signatures to check
var (
    transferSig     = "transfer(address,uint256)"
    transferFromSig = "transferFrom(address,address,uint256)"
)

type TaxDetector struct {
    client *ethclient.Client
	tokenEventTracker *transferEventGather.TokenEventTracker
    myTaxDetectorHelper *MyTaxDetectorHelper
	batchSize         int64 // Number of blocks to analyze per batch
    minTransfers      int   // Minimum transfers to analyze
}

func NewTaxDetector() (*TaxDetector, error) {
	client := blockchainclient.GetHttpClient()

	tokenEventTracker, err := transferEventGather.NewTokenEventTracker(client)
	if err != nil {
		return nil, err
	}

	myTaxDetectorHelper, err := NewMyTaxDetectorHelper()
	if err != nil {
		return nil, err
	}

    return &TaxDetector{
        client: client,
		tokenEventTracker: tokenEventTracker,
		myTaxDetectorHelper: myTaxDetectorHelper,
		batchSize:        1000,  // Analyze 10000 blocks at a time
        minTransfers:     50,     // Minimum transfers to analyze
    }, nil
}

func (td *TaxDetector) Run()  {

	amountOfTokensProcessed := 0
	
	erc20List := erc20.GetAllNonProcessedIsTaxToken()

	for _, element := range erc20List {

		tokenAddress := element.ContractAddress
		
		isTaxToken, err := td.DetectTaxToken(context.Background(), tokenAddress)
		if err == nil {
            erc20.UpdateTax(tokenAddress, isTaxToken, -1)
            pair.UpdateAllPairsWithTokenToAPairWithATaxToken(tokenAddress, isTaxToken)
        } else {
			log.Error().Msgf("Error detecting tax token: %v", err)
		}

		

		amountOfTokensProcessed++

		if amountOfTokensProcessed % 100 == 0 {
			log.Info().Msgf("Processed %d pairs out of %d", amountOfTokensProcessed, len(erc20List))
		}
	}
}


// DetectTaxToken checks if a token has transfer fees
func (td *TaxDetector) DetectTaxToken(ctx context.Context, tokenAddress common.Address) (bool, error) {

    isTax, err := td.checkContractCode(ctx, tokenAddress)
    if err != nil {
        return true, err
    }

    return isTax, nil
}

func (td *TaxDetector) checkTransferEvents(ctx context.Context, tokenAddress common.Address) (bool, float64, error) {
    
	log.Info().Msgf("Checking transfer events for token: %s", tokenAddress.Hex())

	currentBlock := blockchainutil.GetCurrentBlockNumber()
	
	fromBlock := big.NewInt(int64(currentBlock - 100000)) // Look back 100000 blocks
    if fromBlock.Int64() < 0 {
        fromBlock = big.NewInt(0)
    }
    toBlock := new(big.Int).SetUint64(currentBlock)
    
    var allTransfers []transferEventGather.TransferEvent
    taxCount := 0
	
    // Collect transfers in batches
    for fromBlock.Cmp(toBlock) < 0 {
        batchEnd := new(big.Int).Add(fromBlock, big.NewInt(td.batchSize))
        if batchEnd.Cmp(toBlock) > 0 {
            batchEnd = toBlock
        }

        transfers, err := td.tokenEventTracker.GetTransferEvents(ctx, tokenAddress, fromBlock, batchEnd)
        if err != nil {
            return false, 0, err
        }

		transfersWithOutContractAddresses := []transferEventGather.TransferEvent{}

		for _, transfer := range transfers {
			
			log.Debug().Msgf("Transfer: %s", transfer.TxHash.String())
			isContractFrom, err := blockchainutil.IsContractAddress(transfer.From)
			if err != nil {
				log.Error().Msgf("Error checking if address is contract: %v", err)
			}

			if !isContractFrom {

				isContractTo, err := blockchainutil.IsContractAddress(transfer.To)
				if err != nil {
					log.Error().Msgf("Error checking if address is contract: %v", err)
				}

				if !isContractTo{
					transfersWithOutContractAddresses = append(transfersWithOutContractAddresses, transfer)
				}
			}
			
		}

        allTransfers = append(allTransfers, transfersWithOutContractAddresses...)
        
        // If we have enough transfers, start analysis
        if len(allTransfers) >= td.minTransfers {
            break
        }

        fromBlock = batchEnd
    }

    if len(allTransfers) < td.minTransfers {
        return false, 0, nil // Not enough transfers to make a determination
    }

    // Analyze transfers for tax patterns
    totalTransfers := 0
    totalTaxPercentage := 0.0

    for i := 0; i < len(allTransfers)-1; i++ {
        // Look for consecutive transfers from the same transaction
		log.Info().Msgf("Transfer: %s", allTransfers[i].TxHash.String())
        if allTransfers[i].TxHash == allTransfers[i+1].TxHash {
            sent := allTransfers[i].Value
            received := allTransfers[i+1].Value

            // If sent amount is greater than received amount, might be a tax
            if sent.Cmp(received) > 0 {
                diff := new(big.Int).Sub(sent, received)
                percentage := new(big.Float).Quo(
                    new(big.Float).SetInt(diff),
                    new(big.Float).SetInt(sent),
                )
                
                taxPercentage, _ := percentage.Float64()
                taxPercentage *= 100

                // Consider it a tax if the difference is more than 0.1%
                if taxPercentage > 0.1 {
                    taxCount++
                    totalTaxPercentage += taxPercentage
                }
            }
        }
        totalTransfers++
    }

    // Calculate tax frequency and average tax percentage
    if totalTransfers > 0 {
        taxFrequency := float64(taxCount) / float64(totalTransfers)
        averageTaxPercentage := 0.0
        if taxCount > 0 {
            averageTaxPercentage = totalTaxPercentage / float64(taxCount)
        }

        // Consider it a tax token if more than 1% of transfers show tax behavior
        isTaxToken := taxFrequency > 0.01
        return isTaxToken, averageTaxPercentage, nil
    }

    return false, 0, nil
}

// testTransfer performs a test transfer to detect fees
func (td *TaxDetector) testTransfer(ctx context.Context, tokenAddress common.Address) (bool, float64, error) {
    // Create test addresses
    sender := common.HexToAddress("0x1234...")    // Replace with test wallet
    receiver := common.HexToAddress("0x5678...")  // Replace with test wallet

    // Get initial balances
    initialSenderBalance, err := td.getTokenBalance(ctx, tokenAddress, sender)
    if err != nil {
        return false, 0, err
    }

    initialReceiverBalance, err := td.getTokenBalance(ctx, tokenAddress, receiver)
    if err != nil {
        return false, 0, err
    }

    // Perform test transfer
    amount := big.NewInt(1000000) // Use appropriate test amount
    
    // Simulate transfer
    msg := ethereum.CallMsg{
        From:  sender,
        To:    &tokenAddress,
        Value: big.NewInt(0),
        Data:  td.createTransferData(receiver, amount),
    }
    
    _, err = td.client.CallContract(ctx, msg, nil)
    if err != nil {
        return false, 0, err
    }

    // Get final balances
    finalSenderBalance, err := td.getTokenBalance(ctx, tokenAddress, sender)
    if err != nil {
        return false, 0, err
    }

    finalReceiverBalance, err := td.getTokenBalance(ctx, tokenAddress, receiver)
    if err != nil {
        return false, 0, err
    }

    // Calculate actual received amount
    senderDiff := new(big.Int).Sub(initialSenderBalance, finalSenderBalance)
    receiverDiff := new(big.Int).Sub(finalReceiverBalance, initialReceiverBalance)

    // If the received amount is less than sent amount, it's a tax token
    if senderDiff.Cmp(receiverDiff) != 0 {
        taxAmount := new(big.Int).Sub(senderDiff, receiverDiff)
        taxPercentage := float64(taxAmount.Int64()) / float64(senderDiff.Int64()) * 100
        return true, taxPercentage, nil
    }

    return false, 0, nil
}

// checkContractCode analyzes contract bytecode for tax patterns
func (td *TaxDetector) checkContractCode(ctx context.Context, tokenAddress common.Address) (bool, error) {
    code, err := td.client.CodeAt(ctx, tokenAddress, nil)
    if err != nil {
        return false, err
    }

    log.Info().Msgf("Contract Address: %s", tokenAddress.String())
    
    // Convert bytecode to string for pattern matching
    codeStr := "0x" + hex.EncodeToString(code)

    isTaxToken, err := td.myTaxDetectorHelper.IsTaxToken(codeStr)
    if err != nil {
        return true, err
    }

    if isTaxToken {
        log.Info().Msgf("Is tax token: %s", tokenAddress.String())
    }

    return isTaxToken, nil
}

// Helper function to get token balance
func (td *TaxDetector) getTokenBalance(ctx context.Context, tokenAddress, account common.Address) (*big.Int, error) {
    data := td.createBalanceOfData(account)
    
    msg := ethereum.CallMsg{
        To:   &tokenAddress,
        Data: data,
    }
    
    result, err := td.client.CallContract(ctx, msg, nil)
    if err != nil {
        return nil, err
    }
    
    return new(big.Int).SetBytes(result), nil
}

// Helper function to create transfer data
func (td *TaxDetector) createTransferData(to common.Address, amount *big.Int) []byte {
    methodID := crypto.Keccak256([]byte(transferSig))[:4]
    paddedAddress := common.LeftPadBytes(to.Bytes(), 32)
    paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
    
    return append(append(methodID, paddedAddress...), paddedAmount...)
}

// Helper function to create balanceOf data
func (td *TaxDetector) createBalanceOfData(account common.Address) []byte {
    methodID := crypto.Keccak256([]byte("balanceOf(address)"))[:4]
    paddedAddress := common.LeftPadBytes(account.Bytes(), 32)
    
    return append(methodID, paddedAddress...)
}