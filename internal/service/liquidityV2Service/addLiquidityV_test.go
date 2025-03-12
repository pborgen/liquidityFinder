package addLiquidityV2

import (
	"testing"

	"github.com/pborgen/liquidityFinder/myConst"
	"github.com/rs/zerolog/log"
)
func TestBla(t *testing.T) {
	

	log.Info().Msgf("one ddddpls: %s", myConst.GetOneWplsBigint().String())
}

func TestAddLiquidityV2(t *testing.T) {

// 	testhelper.Setup()
	
// 	account, err := account.GetById(14)
// 	if err != nil {
// 		t.Fatalf("Error getting account: %v", err)
// 	}

// 	dexModel := dex.GetById(3)


// 	// erc20ContractAddress, erc20TransactionHash, err := erc20Service.CreateERC20(
// 	// 	account.PrivateKey,
// 	// 	"TeslaTestCoin",
// 	// 	"TSLA",
// 	// 	"TSLA")

// 	if err != nil {
// 		panic(err)
// 	}
// 	erc20ContractAddress := common.HexToAddress("0x0e6f8aC8739e14E2fD11ecEAE60f7437b759D36F")

// //0x0e6f8aC8739e14E2fD11ecEAE60f7437b759D36F
// 	log.Info().Msgf("ERC20 contract address: %s", erc20ContractAddress.Hex())
// 	//log.Info().Msgf("ERC20 transaction hash: %s", erc20TransactionHash)
	

// 	log.Info().Msgf("Router contract address: %s", dexModel.RouterContractAddress.Hex())

	
	
// 	AddInitialLiquidity(
// 		nil,
// 		dexModel.RouterContractAddress,
// 		myConst.GetWplsAddress(),
// 		erc20ContractAddress,
// 		myConst.GetOneThousandWplsBigint(),
// 		myConst.GetOneWplsBigint(),
// 		common.HexToAddress("0x0e6f8aC8739e14E2fD11ecEAE60f7437b759D36F"),
// 	)

}

