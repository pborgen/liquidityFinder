package blockchainclient

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pborgen/liquidityFinder/internal/myConfig"
	"github.com/rs/zerolog/log"
)

var httpClientInstance *ethclient.Client
var webSocketClientInstance *ethclient.Client
var onceGetHttpClient sync.Once
var onceGetWebSocketClient sync.Once
const abiRevert = `[{ "name": "Error", "type": "function", "inputs": [ { "type": "string" } ] }]`


func GetHttpClient() *ethclient.Client {
	onceGetHttpClient.Do(func() {
		
		url := myConfig.GetInstance().BlockchainClientUrlHttp
		
		client, err := ethclient.Dial(url)

		if err != nil {
			log.Fatal().Msgf("Error in GetClient")
		} else {
			log.Log().Msgf("Success! you are connected to the Network")
		}

		httpClientInstance = client
	})
	return httpClientInstance
}

func GetPublicHttpClient() *ethclient.Client {
	onceGetHttpClient.Do(func() {
		
		url := myConfig.GetInstance().BlockchainClientPublicUrlHttp
		
		client, err := ethclient.Dial(url)

		if err != nil {
			log.Fatal().Msgf("Error in GetClient")
		} else {
			log.Log().Msgf("Success! you are connected to the Network")
		}

		httpClientInstance = client
	})
	return httpClientInstance
}

func GetRevertReason(from common.Address, tx *types.Transaction) (string, error) {

	client := GetHttpClient()

	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	res, err := client.CallContract(context.Background(), msg, nil)

	if err != nil {
		if strings.Contains(err.Error(), "reverted") {
			return err.Error(), nil
		}
		return "", err
	}

	return string(res), nil
}

func parseRevertReason(input []byte) (string, error) {
	if len(input) < 4 {
		return "", fmt.Errorf("invalid input")
	}

	// methodID := input[:4]
	inputData := input[4:]

	parsedRevert, err := abi.JSON(strings.NewReader(abiRevert))
	if err != nil {
		return "", err
	}

	var reason string
	err = parsedRevert.UnpackIntoInterface(&reason, "Error", inputData)
	if err != nil {
		return "", err
	}

	return reason, nil
}


func GetWebSocketClient() *ethclient.Client {
	onceGetWebSocketClient.Do(func() {
		
		url := myConfig.GetInstance().BlockchainClientUrlWs
		client, err := ethclient.Dial(url)

		if err != nil {
			log.Fatal().Msgf("Error in GetWebSocketClient")
		} else {
			log.Log().Msgf("Success! you are connected to the Network")
		}

		webSocketClientInstance = client
	})
	return webSocketClientInstance
}


func GetPublicRpcHttpUrl() string {
	return "https://rpc.pulsechain.com"
}
