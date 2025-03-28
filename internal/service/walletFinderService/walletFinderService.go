package walletFinderService

import (
	"log"

	"github.com/pborgen/liquidityFinder/internal/blockchain/blockchainutil"
	"github.com/pborgen/liquidityFinder/internal/database/model/account"
)


func Start() {
	account, err := account.GetById(77)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := blockchainutil.GetPrivateKeyFromAccount(account)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(privateKey)
	log.Println(privateKey)
}