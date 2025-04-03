package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/database/model/dex"
	"github.com/pborgen/liquidityFinder/internal/database/model/erc20"
	"github.com/rs/zerolog/log"
)



type ModelPair struct {
	PairId              	int            		`postgres.Table:"PAIR_ID" json:"pair_id"`
	DexId               	int            		`postgres.Table:"DEX_ID" json:"dex_id"`
	PairIndex           	int            		`postgres.Table:"PAIR_INDEX" json:"pair_index"`
	PairContractAddress 	common.Address 		`postgres.Table:"PAIR_CONTRACT_ADDRESS" json:"pair_contract_address"`
	Token0Erc20Id       	int            		`postgres.Table:"TOKEN0_ID" json:"token0_erc20_id"`
	Token1Erc20Id       	int           		`postgres.Table:"TOKEN1_ID" json:"token1_erc20_id"`
	Token0Address      		common.Address 		`postgres.Table:"TOKEN0_ADDRESS" json:"token0_address"`
	Token1Address      		common.Address 		`postgres.Table:"TOKEN1_ADDRESS" json:"token1_address"`
	Token0Erc20         	erc20.ModelERC20	`json:"token0_erc20"`
	Token1Erc20         	erc20.ModelERC20	`json:"token1_erc20"`
	ModelDex            	dex.ModelDex		`json:"model_dex"`
	Token0Reserves      	big.Int 			`postgres.Table:"TOKEN0_RESERVES" json:"token0_reserves"`
	Token1Reserves      	big.Int 			`postgres.Table:"TOKEN1_RESERVES" json:"token1_reserves"`
	ShouldFindArb       	bool    			`postgres.Table:"SHOULD_FIND_ARB" json:"should_find_arb"`
	IsPlsPair           	bool    			`postgres.Table:"IS_PLS_PAIR" json:"is_pls_pair"`
	HasTaxToken         	bool    			`postgres.Table:"HAS_TAX_TOKEN" json:"has_tax_token"`
	IsHighLiquidity     	bool    			`postgres.Table:"IS_HIGH_LIQUIDITY" json:"is_high_liquidity"`
	UniswapV3Fee           	big.Int            	`postgres.Table:"UNISWAP_V3_FEE"`
	UniswapV3TickSpacings  	big.Int            	`postgres.Table:"UNISWAP_V3_TICK_SPACINGS"`
	InsertedAt          	time.Time 			`postgres.Table:"INSERTED_AT" json:"inserted_at"`
	LastUpdated         	time.Time 			`postgres.Table:"LAST_UPDATED" json:"last_updated"`
	LastTimeReservesUpdated time.Time 			`postgres.Table:"LAST_TIME_RESERVES_UPDATED" json:"last_time_reserves_updated"`
}

type NameValue struct {
	Id int `postgres.Table:"ID"`
	Name string `postgres.Table:"NAME"`
	Value string `postgres.Table:"VALUE"`
	DataType int `postgres.Table:"DATA_TYPE"`
}

type TransferEventGroupBy struct {
	FromAddress common.Address `postgres.Table:"FROM_ADDRESS"`
	ToAddress common.Address `postgres.Table:"TO_ADDRESS"`
	TransactionCount int `postgres.Table:"TRANSACTION_COUNT"`
}
type ModelTransferEvent struct {
	BlockNumber uint64 `postgres.Table:"BLOCK_NUMBER"`
	LogIndex int `postgres.Table:"LOG_INDEX"`
	ContractAddress common.Address `postgres.Table:"CONTRACT_ADDRESS"`
	FromAddress common.Address `postgres.Table:"FROM_ADDRESS"`
	ToAddress common.Address `postgres.Table:"TO_ADDRESS"`
	EventValue *big.Int `postgres.Table:"EVENT_VALUE"`
}

type ModelTokenAmount struct {
	TokenAddress common.Address `postgres.Table:"TOKEN_ADDRESS"`
	OwnerAddress common.Address `postgres.Table:"OWNER_ADDRESS"`
	Amount *big.Int `postgres.Table:"AMOUNT"`
	LastBlockNumberUpdated uint64 `postgres.Table:"LAST_BLOCK_NUMBER_UPDATED"`
	LastLogIndexUpdated int `postgres.Table:"LAST_LOG_INDEX_UPDATED"`
}

func PrintModelTokenAmount(modelTokenAmount *ModelTokenAmount) {
	log.Info().Msgf("modelTokenAmount: %s, %s, %s, %d, %d", modelTokenAmount.TokenAddress.Hex(), modelTokenAmount.OwnerAddress.Hex(), modelTokenAmount.Amount.String(), modelTokenAmount.LastBlockNumberUpdated, modelTokenAmount.LastLogIndexUpdated)
}

type PairsOrganized struct {
	plsPairs    []ModelPair
	nonPlsPairs []ModelPair
}

// END MODEL PAIR


func IsV2Pair(modelPair *ModelPair) bool {
	modelDex := dex.GetById(modelPair.DexId)

	return modelDex.DexType == dex.UniswapV2	
}

func IsV3Pair(modelPair *ModelPair) bool {
	modelDex := dex.GetById(modelPair.DexId)

	return modelDex.DexType == dex.UniswapV3
}