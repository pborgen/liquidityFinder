package accountService

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pborgen/liquidityFinder/internal/database/model/account"
	blockchainutil "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainutil"
	erc20Helper "github.com/pborgen/liquidityFinder/internal/blockchain/erc20"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

func GetById(id int) (*account.ModelAccount, error) {
	return account.GetById(id)
}

func GenerateAccount(name string, description string, isActive bool) (*account.ModelAccount, error) {

	publicKey, privateKey := blockchainutil.GeneratePublicAndPrivateKey()

	return account.Insert(name, description, publicKey, privateKey, isActive)
}

func AddAccount(name string, description string, isActive bool, publicKey string, privateKey string) (*account.ModelAccount, error) {


	account, err := account.Insert(name, description, publicKey, privateKey, isActive)
	if err != nil {
		return nil, err
	}

	return account, nil
}
func GeneratePublicKeyStartingWith(startingWithArray []string) (results[]struct {
	PublicKey  string
	PrivateKey string
}, err error) {


	wg := sync.WaitGroup{}
	sem := semaphore.NewWeighted(1000)		
	resultChan := make(chan struct {
		publicKey  string
		privateKey string
	}, 1)

	for i := 0; i < 250_000_000; i++ {
		wg.Add(1)
		sem.Acquire(context.Background(), 1)

		go func() {
			defer wg.Done()
			defer sem.Release(1)
			publicKey, privateKey := blockchainutil.GeneratePublicAndPrivateKey()

			for _, startingWith := range startingWithArray {
				if strings.HasPrefix(publicKey, startingWith) {
					
					resultChan <- struct {
						publicKey  string
						privateKey string
					}{publicKey, privateKey}
					return
				}
			}
		}()
	}
	
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for element := range resultChan {
		results = append(results, struct {
			PublicKey  string
			PrivateKey string
		}{element.publicKey, element.privateKey})
	}	

	if len(results) > 0 {
		return results, nil
	}

	return nil, errors.New("could not generate a public key starting with any of " + strings.Join(startingWithArray, ", "))
}

// func MovePls(account *account.ModelAccount, toAddress common.Address, amount *big.Int) (bool, error) {
// 	auth, err := blockchainutil.GetAuthAccount(account.PrivateKey, 369)
// 	if err != nil {
// 		return false, err
// 	}

// 	success, err := blockchainutil.SendPls(auth, account, toAddress, amount)
// 	if err != nil {
// 		return false, err
// 	}

// 	return success, nil
// }

func MoveTokensInAccount(
	account *account.ModelAccount, 
	tokenAddress string, 
	toAddress string, 
	maxGasPriceInPls *big.Int,
) (tx *types.Transaction, err error) {

	auth, err := blockchainutil.GetAuthAccount(account.PrivateKey, 369)
	if err != nil {
		return nil, err
	}

	balance, err := erc20Helper.BalanceOf( common.HexToAddress(tokenAddress), common.HexToAddress(account.PublicKey))
	if err != nil {
		return nil, err
	}

	if balance.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("no balance")
	} else {
		plsCost, err := 
			erc20Helper.EstimatePLSNeededForTransfer(
				auth, 
				common.HexToAddress(tokenAddress), 
				common.HexToAddress(toAddress), 
				balance,
			)
		if err != nil {
			return nil, err
		}

		log.Info().Msgf("Estimated PLS cost for transfer: %s", plsCost.String())

		// if plsCost.Cmp(maxGasPriceInPls) > 0 {
		// 	return errors.New("not enough PLS to transfer")
		// }

		tx, err := 
			erc20Helper.Transfer(
				auth, 
				common.HexToAddress(tokenAddress), 
				common.HexToAddress(toAddress), 
				balance,
				uint64(60000),
			)
		if err != nil {
			return nil, err
		}
	
		return tx, nil
	}

	return nil, errors.New("no balance")
}

