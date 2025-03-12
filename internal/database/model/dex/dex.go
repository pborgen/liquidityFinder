package dex

import (
	"database/sql"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/database"
	"github.com/pborgen/liquidityFinder/internal/database/model/orm"
	"github.com/rs/zerolog/log"
)

const (
	UniswapV2 int = 1
	UniswapV3 int = 2
)

type ModelDex struct {
	DexId                  int            `postgres.Table:"DEX_ID"`
	Name                   string         `postgres.Table:"NAME"`
	NetworkId              int            `postgres.Table:"NETWORK_ID"`
	RouterContractAddress  common.Address `postgres.Table:"ROUTER_ADDRESS"`
	FactoryContractAddress common.Address `postgres.Table:"FACTORY_ADDRESS"`
	RouterAbi              string         `postgres.Table:"ROUTER_ABI"`
	FactoryAbi             string         `postgres.Table:"FACTORY_ABI"`
	FeeBasisPoints         int64          `postgres.Table:"FEE_BASIS_POINTS"`
	DexType                int            `postgres.Table:"DEX_TYPE"`
}

var dexColumnNames = orm.GetColumnNames(ModelDex{})

const tableName = "DEX"

func GetAllV2Dexes(networkId int) []ModelDex {

	allDexes := GetAll()

	filteredDexes := make([]ModelDex, 0)
	for _, dex := range allDexes {
		if dex.DexType == UniswapV2 && dex.NetworkId == networkId {
			filteredDexes = append(filteredDexes, dex)
		}
	}

	return filteredDexes
}

func GetAllV3Dexes(networkId int) []ModelDex {

	allDexes := GetAll()
	
	filteredDexes := make([]ModelDex, 0)
	for _, dex := range allDexes {
		if dex.DexType == UniswapV3 && dex.NetworkId == networkId {
			filteredDexes = append(filteredDexes, dex)
		}
	}

	return filteredDexes
}

func GetById(id int) ModelDex {

	db := database.GetDBConnection()
	query := "SELECT " + dexColumnNames + " FROM " + tableName + " WHERE DEX_ID = $1"
	rows, err := db.Query(query, id)

	if err != nil {
		panic("Error getting DexModel with id:" + strconv.Itoa(id))
	}

	defer rows.Close()

	if rows.Next() {
		dex, err := scan(rows)

		if err != nil {
			panic("Error getting DexModel with id:" + strconv.Itoa(id))
		}
		return *dex
	} else {
		panic("Could not find dex with id:" + strconv.Itoa(id))
	}

}

func GetAllDexIdsByNetworkId(networkId int) []int {
	allDexes := GetAll()

	filteredDexes := make([]int, 0)
	for _, dex := range allDexes {
		if dex.NetworkId == networkId {
			filteredDexes = append(filteredDexes, dex.DexId)
		}
	}
	return filteredDexes
}

func GetAllByNetworkId(networkId int) []ModelDex {
	allDexes := GetAll()

	filteredDexes := make([]ModelDex, 0)
	for _, dex := range allDexes {
		if dex.NetworkId == networkId {
			filteredDexes = append(filteredDexes, dex)
		}
	}
	return filteredDexes
}

func ModelDexToMap(modelDexList[]ModelDex) (map[int]ModelDex) {
	modelDexMap := make(map[int]ModelDex)

	for _, modelDex := range modelDexList {
		modelDexMap[modelDex.DexId] = modelDex
	}

	return modelDexMap
}

func ModelDexToMapWithRouterAddressAsKey(modelDexList[]ModelDex) (map[common.Address]ModelDex) {
	modelDexMap := make(map[common.Address]ModelDex)

	for _, modelDex := range modelDexList {
		modelDexMap[modelDex.RouterContractAddress] = modelDex
	}

	return modelDexMap
}

func GetAll() []ModelDex {

	db := database.GetDBConnection()

	results := make([]ModelDex, 0)
	rows, err := db.Query("SELECT " + dexColumnNames + " FROM " + tableName)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		dex, err := scan(rows)
		if err != nil {

			panic(err)
		}
		results = append(results, *dex)
	}

	return results
}

func scan(rows *sql.Rows) (*ModelDex, error) {
	
	dex := ModelDex{}

	var routerContractAddressString string
	var factoryContractAddressString string
	var routerAbiString sql.NullString
	var factoryAbiString sql.NullString

	err := rows.Scan(
		&dex.DexId, 
		&dex.Name, 
		&dex.NetworkId, 
		&routerContractAddressString, 
		&factoryContractAddressString, 
		&routerAbiString, 
		&factoryAbiString, 
		&dex.FeeBasisPoints,
		&dex.DexType,
	)
	if err != nil {
		log.Error().Msgf("Could not scan rows for dex", err)

	}

	dex.RouterContractAddress = common.HexToAddress(routerContractAddressString)
	dex.FactoryContractAddress = common.HexToAddress(factoryContractAddressString)

	dex.RouterAbi = routerAbiString.String
	dex.FactoryAbi = factoryAbiString.String

	return &dex, nil
}
