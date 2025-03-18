package updateSymbol

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/myConfig"
	"github.com/pborgen/liquidityFinder/internal/myJsonHelper"
	"github.com/pborgen/liquidityFinder/internal/mylogger"
	"github.com/pborgen/liquidityFinder/internal/service/erc20Service"
	"github.com/rs/zerolog/log"
)


func Start() {
	mylogger.Init()

	instance := myConfig.GetInstance()
	baseDir := instance.BaseDir


	jsonData, err := myJsonHelper.ReadJSON(baseDir + "/cmd/updateSymbol/symbolImages.json")
	if err != nil {
		
		panic(err)
	}

	for _, token := range jsonData["tokens"].([]interface{}) {
		tokenMap := token.(map[string]interface{})
		logoURI := tokenMap["logoURI"].(string)
		address := tokenMap["address"].(string)

		log.Info().Str("logoURI", logoURI).Str("address", address).Msg("symbol")

		erc20Service.UpdateSymbolImageUrl(common.HexToAddress(address), logoURI)
	}

	log.Info().Interface("jsonData", jsonData).Msg("jsonData")
	
}

