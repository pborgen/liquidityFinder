package blockchainclient

import (
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

var httpClientInstance *ethclient.Client
var webSocketClientInstance *ethclient.Client
var onceGetHttpClient sync.Once
var onceGetWebSocketClient sync.Once


func GetHttpClient() *ethclient.Client {
	onceGetHttpClient.Do(func() {
		
		url := "https://rpc.pulsechain.com"
		
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

func GetWebSocketClient() *ethclient.Client {
	onceGetWebSocketClient.Do(func() {
		
		url := "wss://rpc.pulsechain.com"
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
