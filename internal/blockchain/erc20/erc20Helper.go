package erc20Helper

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pborgen/liquidityFinder/abi/erc20"
	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
	blockchainutil "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainutil"
	"github.com/rs/zerolog/log"
)

func GetErc20(contractAddress common.Address) (*struct {
	Name    string
	Symbol  string
	Decimal uint8
}, error) {

	log.Info().Msgf("Getting Erc20 with address:" + contractAddress.String())
	client := blockchainclient.GetHttpClient()

	contract, err := erc20.NewErc20(contractAddress, client)

	if err != nil {
		return nil, fmt.Errorf("could not create the contract: %w", err)
	}

	name, err := contract.Name(nil)

	if err != nil {

		return nil, fmt.Errorf("could not get the name of the token: %w", err)
	}

	symbol, err := contract.Symbol(nil)

	if err != nil {
		return nil, fmt.Errorf("could not get the symbol of the token: %w", err)
	}

	decimal, err := contract.Decimals(nil)

	if err != nil {
		return nil, fmt.Errorf("could not get the decimals of the token: %w", err)
	}

	outStruct := new(struct {
		Name    string
		Symbol  string
		Decimal uint8
	})


	outStruct.Name = name
	outStruct.Symbol = symbol
	outStruct.Decimal = decimal


	return outStruct, err
}

func IsApproved(
	tokenAddress common.Address,
	owner common.Address,
	spender common.Address,
	amount *big.Int,
) (bool, error) {
	client := blockchainclient.GetHttpClient()

	contract, err := erc20.NewErc20Caller(tokenAddress, client)
	if err != nil {
		return false, err
	}

	allowance, err := contract.Allowance(nil, owner, spender)
	if err != nil {
		return false, err
	}

	if allowance.Cmp(amount) >= 0 {
		return true, nil
	}

	return false, nil
}

func Approve(
	auth *bind.TransactOpts,
	tokenAddress common.Address,
	spender common.Address,
	amount *big.Int,
) (bool, error) {
	client := blockchainclient.GetHttpClient()

	log.Info().Msgf("Approving token: %s with spender: %s and amount: %s", tokenAddress.Hex(), spender.Hex(), amount.String())

	contract, err := erc20.NewErc20Transactor(tokenAddress, client)
	if err != nil {
		return false, err
	}

	approved, err := IsApproved(tokenAddress, auth.From, spender, amount)
	if err != nil {
		return false, err
	}

	if approved {
		log.Info().Msgf("Token already approved. TokenAddress: %s, spender: %s, amount: %s", tokenAddress.String(), spender.String(), amount.String())
		return true, nil
	}

	var success bool

	if !approved {

		err = blockchainutil.SetupAuthForTransaction(auth, false)
		if err != nil {
			return false, err
		}

		transaction, err := contract.Approve(auth, spender, amount)

		if err != nil {
			return false, err
		}
		time.Sleep(5 * time.Second)

		log.Info().Msgf(
			"Approving: transaction: %s For tokenAddress: %s and spender: %s", 
			transaction.Hash().String(), 
			tokenAddress.String(), 
			spender.String(),
		)

		success = true

		// success, err = blockchainutil.GetTransactionResults(auth, transaction)
		// if err != nil {
		// 	return false, err
		// }
	
		// log.Info().Msgf(
		// 	"Approved: transaction: %s For tokenAddress: %s and spender: %s", 
		// 	transaction.Hash().String(), 
		// 	tokenAddress.String(), 
		// 	spender.String(),
		// )
	}

	return success, nil
}

func BalanceOf(
	tokenAddress common.Address,
	owner common.Address,
) (*big.Int, error) {
	client := blockchainclient.GetHttpClient()

	contract, err := erc20.NewErc20Caller(tokenAddress, client)
	if err != nil {
		return nil, err
	}

	balance, err := contract.BalanceOf(nil, owner)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func Transfer(
	auth *bind.TransactOpts,
	tokenAddress common.Address,
	toAddress common.Address,
	amount *big.Int,
	gasLimit uint64, // 0 = estimate
) (tx *types.Transaction, err error) {
	client := blockchainclient.GetHttpClient()

	blockchainutil.SetupAuthForTransaction(auth, false)

	if gasLimit != 0 {	
		auth.GasLimit = gasLimit
	}

	contract, err := erc20.NewErc20Transactor(tokenAddress, client)
	if err != nil {
		return nil, err
	}

	transaction, err := contract.Transfer(auth, toAddress, amount)
	if err != nil {
		return nil, err
	}
	

	log.Info().Msgf("Transfer tx: %s", transaction.Hash().Hex())

	return transaction, nil
}

func EstimatePLSNeededForTransfer(
	auth *bind.TransactOpts,
	tokenAddress common.Address,
	toAddress common.Address,
	amount *big.Int,
) (*big.Int, error) {
	client := blockchainclient.GetHttpClient()


	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.Keccak256Hash(transferFnSignature)
	methodID := hash[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// Estimate gas for the transfer
	msg := ethereum.CallMsg{
		From:  auth.From,
		To:    &tokenAddress,
		Value: big.NewInt(0), // No ETH is being transferred, just tokens
		Data:  data,
	}

	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		return nil, err
	}

	// Get the current gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// Calculate the total cost in wei
	totalCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)

	// Convert wei to PLS (assuming 1 PLS = 10^18 wei)
	plsCostFloat := new(big.Float).Quo(new(big.Float).SetInt(totalCost), big.NewFloat(1e18))

	plsCostInt := new(big.Int)
	plsCostFloat.Int(plsCostInt)

	return plsCostInt, nil
}
