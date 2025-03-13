package pair

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/database"
	"github.com/pborgen/liquidityFinder/internal/database/model/dex"
	"github.com/pborgen/liquidityFinder/internal/database/model/erc20"
	"github.com/pborgen/liquidityFinder/internal/database/model/orm"
	"github.com/pborgen/liquidityFinder/internal/myUtil"
	cacheService "github.com/pborgen/liquidityFinder/internal/service/CacheService"
	"github.com/pborgen/liquidityFinder/internal/types"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)


const primaryKey = "PAIR_ID"
const tableName = "PAIR"

var wplsAddress = common.HexToAddress("0xA1077a294dDE1B09bB078844df40758a5D0f9a27")

var pairColumnNames = orm.GetColumnNames(types.ModelPair{})
var wplsModelErc20 erc20.ModelERC20
var err error

func init() {
	wplsModelErc20, err = erc20.GetByContractAddress(wplsAddress, 1)
	if err != nil {
		panic(err)
	}
}

func Insert(pairModel types.ModelPair) (*types.ModelPair, error) {
	db := database.GetDBConnection()

	// Check if ERC20's Exists
	erc20.GetById(pairModel.Token0Erc20Id)
	erc20.GetById(pairModel.Token1Erc20Id)

	sqlStatement := orm.CreateInsertStatement(types.ModelPair{}, tableName, primaryKey)

	id := 0

	err := db.QueryRow(
		sqlStatement, 
		pairModel.DexId, 
		pairModel.PairIndex, 
		pairModel.PairContractAddress.String(), 
		pairModel.Token0Erc20Id, 
		pairModel.Token1Erc20Id, 
		pairModel.Token0Address.String(), 
		pairModel.Token1Address.String(), 
		pairModel.Token0Reserves.String(), 
		pairModel.Token1Reserves.String(), 
		pairModel.ShouldFindArb, 
		pairModel.HasTaxToken,
		pairModel.IsPlsPair, 
		pairModel.IsHighLiquidity,
		pairModel.UniswapV3Fee.String(),
		pairModel.UniswapV3TickSpacings.String(),
		pairModel.InsertedAt,
		pairModel.LastUpdated,
		pairModel.LastTimeReservesUpdated,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Pair Inserted with address + " + pairModel.PairContractAddress.String() + " and dex_id: " + strconv.Itoa(pairModel.DexId))
	return GetById(id)
}

// A function that check if a pair exists
func ExistsByContractAddress(contractAddress common.Address) bool {

	exists := false
	count := 0

	db := database.GetDBConnection()
	row := db.QueryRow(
		"SELECT COUNT(1) FROM "+tableName+" WHERE PAIR_CONTRACT_ADDRESS = $1", contractAddress.String())

	row.Scan(&count)

	if count == 1 {
		exists = true
	}

	return exists
}

func UpdateReserves(pairId int, reserve0 big.Int, reserve1 big.Int) error {
	db := database.GetDBConnection()

	sql := "UPDATE " + tableName + " SET TOKEN0_RESERVES = $1, TOKEN1_RESERVES = $2, LAST_UPDATED = NOW(), LAST_TIME_RESERVES_UPDATED = NOW() WHERE PAIR_ID = $3"

	_, err := db.Exec(sql, reserve0.String(), reserve1.String(), pairId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateForFix(modelPair types.ModelPair) error {
	db := database.GetDBConnection()

	sql := "UPDATE " + tableName + " SET TOKEN0_ID=$1, TOKEN1_ID=$2, TOKEN0_ADDRESS=$3, TOKEN1_ADDRESS=$4, TOKEN0_RESERVES=0, TOKEN1_RESERVES=0, LAST_UPDATED = NOW() WHERE PAIR_ID = $5"

	result, err := db.Exec(
		sql, 
		modelPair.Token0Erc20Id, 
		modelPair.Token1Erc20Id, 
		modelPair.Token0Address.String(), 
		modelPair.Token1Address.String(), 
		modelPair.PairId)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("No rows affected")
	}

	return nil
}

func UpdateTax(pairId int, hasTaxToken bool) error {
	db := database.GetDBConnection()

	sql := "UPDATE " + tableName + " SET HAS_TAX_TOKEN = $1, LAST_UPDATED = NOW() WHERE PAIR_ID = $2"

	_, err := db.Exec(sql, hasTaxToken, pairId)
	if err != nil {
		return err
	}

	return nil
}

func GetById(id int) (*types.ModelPair, error) {

	db := database.GetDBConnection()

	rows := db.QueryRow(
		"SELECT "+pairColumnNames+" FROM "+tableName+" WHERE PAIR_ID = $1", id)

	if rows == nil {
		return nil, errors.New("Could not find pairId:" + strconv.Itoa(id))
	}

	pairModel, err := scan(rows, true)

	if err != nil {
		return nil, err
	}

	return pairModel, nil
}

func GetRandomPairs() ([]types.ModelPair, error) {
	db := database.GetDBConnection()
	results := []types.ModelPair{}
	rows, err := db.Query("SELECT " + pairColumnNames + " FROM " + tableName + " ORDER BY RANDOM() LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, true)
		if err != nil {
			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}

func GetPairsByContractAddresses(contractAddresses []common.Address) ([]types.ModelPair, error) {
	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)

	contractAddressesStringList := make([]string, len(contractAddresses))
	for i, contractAddress := range contractAddresses {
		contractAddressesStringList[i] = fmt.Sprintf("'%s'", contractAddress.String())
	}

	query := fmt.Sprintf(
		"SELECT " + pairColumnNames + " FROM " + tableName + " WHERE PAIR_CONTRACT_ADDRESS IN (%s)", 
		strings.Join(contractAddressesStringList, ","),
	)

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, false)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}

func GetByContractAddress(contractAddress common.Address) (*types.ModelPair, error) {

	db := database.GetDBConnection()
	contractAddressAsString := contractAddress.String()
	sql := "SELECT " + pairColumnNames + " FROM " + tableName + " WHERE PAIR_CONTRACT_ADDRESS = $1"
	rows := db.QueryRow(sql, contractAddressAsString)

	pairModel := &types.ModelPair{}

	pairModel, err := scan(rows, true)

	if err != nil {
		return nil, err
	}

	return pairModel, nil
}

func GetSpecificPair(
	dex dex.ModelDex, 
	token0ContractAddress common.Address, 
	token1ContractAddress common.Address) (types.ModelPair, error) {

	db := database.GetDBConnection()

	sql := "SELECT " + pairColumnNames + " FROM " + tableName + " WHERE " +
		"DEX_ID = $1 AND ((TOKEN0_ADDRESS = $2 AND TOKEN1_ADDRESS = $3) OR (TOKEN0_ADDRESS = $4 AND TOKEN1_ADDRESS = $5))"
	
	rows := db.QueryRow(
		sql, 
		dex.DexId, 
		token0ContractAddress.Hex(), 
		token1ContractAddress.Hex(), 
		token1ContractAddress.Hex(), 
		token0ContractAddress.Hex(),
	)

	pairModel := &types.ModelPair{}

	pairModel, err := scan(rows, true)

	if err != nil {
		return types.ModelPair{}, err
	}

	return *pairModel, nil
}

func GetAllWithLimit(limit int, shouldHydrate bool) ([]types.ModelPair, error) {

	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)
	rows, err := db.Query("SELECT " + pairColumnNames + " FROM " + tableName + " LIMIT $1", limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, shouldHydrate)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}

func GetAllPageAndLimit(page int, limit int, shouldHydrate bool) ([]types.ModelPair, error) {

	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)
	rows, err := db.Query("SELECT " + pairColumnNames + " FROM " + tableName + " ORDER BY PAIR_ID ASC LIMIT $1 OFFSET $2", limit, (page-1)*limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, shouldHydrate)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}

func GetAllWithDexIdPageAndLimit(dexId int, page int, limit int, shouldHydrate bool) ([]types.ModelPair, error) {

	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)
	rows, err := db.Query("SELECT " + pairColumnNames + " FROM " + tableName + " WHERE DEX_ID = $1 ORDER BY PAIR_ID ASC LIMIT $2 OFFSET $3", dexId, limit, (page-1)*limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, shouldHydrate)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}

func GetAll(shouldHydrate bool) ([]types.ModelPair, error) {

	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)
	rows, err := db.Query("SELECT " + pairColumnNames + " FROM " + tableName + " ORDER BY PAIR_ID ASC")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, shouldHydrate)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}

func GetAllHighLiquidityPairsThatAreNotTaxToken(shouldHydrate bool) ([]types.ModelPair, error) {
	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)
	rows, err := db.Query("SELECT "+pairColumnNames+" FROM "+tableName+" WHERE IS_HIGH_LIQUIDITY = true AND HAS_TAX_TOKEN = false")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, shouldHydrate)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}

func GetAllNonTaxPairs() ([]types.ModelPair, error) {

	db := database.GetDBConnection()

	cacheServiceInstance := cacheService.GetInstance()
	cacheStruct := cacheService.CacheType_AllNonTaxPairs
	cacheKey := cacheStruct.Name

	// Start Cache
	returnValue, err := cacheService.GetObject[[]types.ModelPair](context.Background(), cacheKey, cacheStruct)

	if err != nil {
		log.Error().Msgf(err.Error())
	}

	if len(returnValue) > 0 {
		return returnValue, nil
	}
	// End Cache
	

	results := make([]types.ModelPair, 0)
	rows, err := db.Query("SELECT "+pairColumnNames+" FROM "+tableName+" WHERE HAS_TAX_TOKEN = false")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, true)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	err = cacheServiceInstance.SetObject(context.Background(), cacheKey, results, cacheStruct)
	if err != nil {
		log.Error().Msgf(err.Error())
	}

	return results, nil

}


func PlsPairWithHighAmountOfPls(dexIds []int, minAmountOfPls *big.Int, shouldHydrate bool) []types.ModelPair {

	var returnValue = []types.ModelPair{}

	cacheServiceInstance := cacheService.GetInstance()
	cacheStruct := cacheService.CacheType_PlsPairWithHighAmountOfPls


	allPairs := []types.ModelPair{}
	for _, dexId := range dexIds {
		
		var allPairsForDexId []types.ModelPair
		cacheKey := cacheStruct.Name + "_" + strconv.Itoa(dexId)

		// Start Cache
		allPairsForDexId, err := cacheService.GetObject[[]types.ModelPair](context.Background(), cacheKey, cacheStruct)
		
		if err != nil {
			log.Error().Msgf(err.Error())
		}

		// End Cache
		if len(allPairsForDexId) == 0 {
			allPairsForDexId, err = GetAllPairsThatHavePlsByDexId(dexId, shouldHydrate)
			if err != nil {
				log.Error().Msgf(err.Error())
			} else {
				if len(allPairsForDexId) > 0 {
					err = cacheServiceInstance.SetObject(context.Background(), cacheKey, allPairsForDexId, cacheStruct)
					if err != nil {
						log.Error().Msgf(err.Error())
					}
				}
			}
		}

		allPairs = append(allPairs, allPairsForDexId...)
	}

	for _, pair := range allPairs {
		shouldAdd := PlsPairHasHighAmountOfPls(pair, minAmountOfPls)

		if shouldAdd {
			returnValue = append(returnValue, pair)
		}
	}

	return returnValue
}

func GetAllPairsThatHavePlsByDexId(dexId int, shouldHydrate bool) ([]types.ModelPair, error) {

	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)

	query := "SELECT "+pairColumnNames+" FROM "+tableName+" WHERE (TOKEN0_ID = $1 OR TOKEN1_ID = $2) AND DEX_ID = $3"
	rows, err := db.Query(query, wplsModelErc20.Erc20Id, wplsModelErc20.Erc20Id, dexId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, shouldHydrate)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}
func GetPairsByPairContractAddresses(pairContractAddresses []common.Address, shouldHydrate bool) ([]types.ModelPair, error) {

	db := database.GetDBConnection()
	pairContractAddressesAsStrings := make([]string, len(pairContractAddresses))
	for i, pairContractAddress := range pairContractAddresses {
		pairContractAddressesAsStrings[i] = pairContractAddress.String()
	}

	query := fmt.Sprintf(
		"SELECT "+pairColumnNames+" FROM "+tableName+" WHERE PAIR_CONTRACT_ADDRESS IN (%s)", 
		strings.Join(pairContractAddressesAsStrings, ","),
	)

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]types.ModelPair, 0)
	for rows.Next() {
		pair, err := scan(rows, shouldHydrate)
		if err != nil {
			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}
func GetAllNonPlsPairsWithHighLiquidity(dexIds []int, shouldHydrate bool) ([]types.ModelPair, error) {

	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)

	placeholders := make([]string, len(dexIds))
	for i := range dexIds {
		placeholders[i] = fmt.Sprintf("%d", dexIds[i])
	}

	query := fmt.Sprintf(`
		SELECT ` +pairColumnNames + ` 
		FROM `+ tableName +` 
		WHERE (TOKEN0_ID != $1 AND TOKEN1_ID != $2) AND 
			  IS_HIGH_LIQUIDITY = true AND DEX_ID IN (%s)`, strings.Join(placeholders, ","))
	
	rows, err := db.Query(query, wplsModelErc20.Erc20Id, wplsModelErc20.Erc20Id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, shouldHydrate)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}

func GetAllPairsThatHavePlsByDexIds(dexIds []int) ([]types.ModelPair, error) {

	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)

	var IDs []string
	for _, i := range dexIds {
		IDs = append(IDs, strconv.Itoa(i))
	}

	query := fmt.Sprintf("SELECT "+pairColumnNames+" FROM "+tableName+" WHERE (TOKEN0_ID = $1 OR TOKEN1_ID = $2) AND DEX_ID IN (%s)", strings.Join(IDs, ","))
	rows, err := db.Query(query, wplsModelErc20.Erc20Id, wplsModelErc20.Erc20Id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, false)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}

func ModelPairToMap(ModelPairList []types.ModelPair) (map[common.Address]types.ModelPair) {
	modelPairMap := make(map[common.Address]types.ModelPair)

	for _, modelPair := range ModelPairList {
		modelPairMap[modelPair.PairContractAddress] = modelPair
	}

	return modelPairMap
}

func GetAllPairsThatHavePls() ([]types.ModelPair, error) {

	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)
	rows, err := db.Query("SELECT "+pairColumnNames+" FROM "+tableName+" WHERE TOKEN0_ID = $1 OR TOKEN1_ID = $1", wplsModelErc20.Erc20Id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, true)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}



func UpdateIsHighLiquidity(pairId int, isHighLiquidity bool) {
	db := database.GetDBConnection()

	query := "UPDATE " + tableName + " SET IS_HIGH_LIQUIDITY = $1, LAST_UPDATED = NOW() WHERE PAIR_ID = $2"
	_, err := db.Exec(query, isHighLiquidity, pairId)

	if err != nil {
		panic(err)
	}


}

func GetAllPairsOnDex(dexId int) ([]types.ModelPair, error) {

	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)
	rows, err := db.Query("SELECT "+pairColumnNames+" FROM "+tableName+" WHERE DEX_ID != $1", dexId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, true)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}

func UpdateAllPairsWithTokenToAPairWithATaxToken(tokenAddress common.Address, isTaxToken bool) (error) {
	db := database.GetDBConnection()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	// Set the transaction isolation level to READ COMMITTED
	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	sql := "UPDATE " + tableName + " SET HAS_TAX_TOKEN = $1, LAST_UPDATED = NOW() WHERE TOKEN0_ADDRESS = $2 OR TOKEN1_ADDRESS = $2"

	_, err = tx.Exec(sql, isTaxToken, tokenAddress.String())
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error updating all pairs with token to a pair with a tax token: %v", err)
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return nil
}

func GetAllPairsWithOutPls() ([]types.ModelPair, error) {

	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)
	rows, err := db.Query("SELECT "+pairColumnNames+" FROM "+tableName+" WHERE TOKEN0_ID != $1 AND TOKEN1_ID != $1", wplsModelErc20.Erc20Id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pair, err := scan(rows, true)

		if err != nil {

			return nil, err
		}
		results = append(results, *pair)
	}

	return results, nil
}

//func GetAllPairsOrganized() ([]types.ModelPair, []types.ModelPair, error) {
//	// Get all the pairs that have pls in it
//	pairs, err := GetAll()
//
//	var plsPairs []types.ModelPair
//	var nonPlsPairs []types.ModelPair
//
//	if err != nil {
//		return nil, nil, err
//	}
//
//	for _, types.ModelPair := range pairs {
//		if types.ModelPair.Token0Erc20.ContractAddress == wplsAddress || modelPair.Token1Erc20.ContractAddress == wplsAddress {
//			plsPairs = append(plsPairs, modelPair)
//		} else {
//			nonPlsPairs = append(nonPlsPairs, modelPair)
//		}
//	}
//
//	return plsPairs, nonPlsPairs, nil
//}

func GetLargestPairIndex(dexModel dex.ModelDex) (int, error) {
	db := database.GetDBConnection()

	var returnValue int = -1
	var count sql.NullInt64

	row := db.QueryRow("SELECT MAX(PAIR_INDEX) FROM "+tableName+" WHERE DEX_ID=$1", dexModel.DexId)
	err := row.Scan(&count)

	if err != nil {
		return -1, err
	}

	if count.Valid {
		returnValue = int(count.Int64)
	}

	return returnValue, nil
}

func GetNonPlsAddress(modelPair types.ModelPair) common.Address {

	return GetNonPlsERC20(modelPair).ContractAddress
}

func GetNonPlsERC20(modelPair types.ModelPair) erc20.ModelERC20 {
	token0ContractAddress := modelPair.Token0Erc20.ContractAddress
	token1ContractAddress := modelPair.Token1Erc20.ContractAddress

	if token0ContractAddress != wplsAddress {
		return modelPair.Token0Erc20
	} else if token1ContractAddress != wplsAddress {
		return modelPair.Token1Erc20
	} else {
		panic("modelPair did not have a wpls address in it. ModlePareId:" + strconv.Itoa(modelPair.PairId))
	}
}

func query(sem *semaphore.Weighted, ch chan []types.ModelPair, wg *sync.WaitGroup, sql string, limit int, offset int) {
	defer wg.Done()
	db := database.GetDBConnection()

	results := make([]types.ModelPair, 0)

	rows, err := db.Query(sql, limit, offset)
	defer rows.Close()

	if err != nil {
		panic(err)
	}

	for rows.Next() {

		pair, err := scan(rows, true)

		if err != nil {

			panic(err)
		}
		results = append(results, *pair)
	}

	log.Log().Msgf("Bock Complete")
	sem.Release(1)
	ch <- results
}

func scan(rows orm.Scannable, shouldHydrate bool) (*types.ModelPair, error) {
	pairModel := types.ModelPair{}

	var tempPairContractAddress string

	var tempToken0Reserves []uint8
	var tempToken1Reserves []uint8
	var tempUniswapV3Fee []uint8
	var tempUniswapV3TickSpacings []uint8

	var tempToken0Address string
	var tempToken1Address string

	err := rows.Scan(
		&pairModel.PairId,
		&pairModel.DexId,
		&pairModel.PairIndex,
		&tempPairContractAddress,
		&pairModel.Token0Erc20Id,
		&pairModel.Token1Erc20Id,
		&tempToken0Address,
		&tempToken1Address,
		&tempToken0Reserves,
		&tempToken1Reserves,
		&pairModel.ShouldFindArb,
		&pairModel.IsPlsPair,
		&pairModel.HasTaxToken,
		&pairModel.IsHighLiquidity,
		&tempUniswapV3Fee,
		&tempUniswapV3TickSpacings,
		&pairModel.InsertedAt,
		&pairModel.LastUpdated,
		&pairModel.LastTimeReservesUpdated,
	)

	if err != nil {
		return &types.ModelPair{}, err
	}

	pairModel.PairContractAddress = common.HexToAddress(tempPairContractAddress)

	tempToken0ReservesString := string(tempToken0Reserves)
	tempToken1ReservesString := string(tempToken1Reserves)
	tempUniswapV3FeeString := string(tempUniswapV3Fee)
	tempUniswapV3TickSpacingsString := string(tempUniswapV3TickSpacings)

	temp0, success := new(big.Int).SetString(tempToken0ReservesString, 10)  // base 10
	if !success {
		log.Error().Msgf("Failed to convert string to big.Int")
	}


	temp1, success := new(big.Int).SetString(tempToken1ReservesString, 10)  // base 10
	if !success {
		log.Error().Msgf("Failed to convert string to big.Int")
	}

	temp2, success := new(big.Int).SetString(tempUniswapV3FeeString, 10)  // base 10
	if !success {
		log.Error().Msgf("Failed to convert string to big.Int")
	}

	temp3, success := new(big.Int).SetString(tempUniswapV3TickSpacingsString, 10)  // base 10
	

	pairModel.Token0Reserves = *temp0
	pairModel.Token1Reserves = *temp1
	pairModel.UniswapV3Fee = *temp2
	pairModel.UniswapV3TickSpacings = *temp3

	pairModel.Token0Address = common.HexToAddress(tempToken0Address)
	pairModel.Token1Address = common.HexToAddress(tempToken1Address)

	if shouldHydrate {
		// Hydrate
		out1 := myUtil.MyAsync(func() erc20.ModelERC20 {
			return erc20.GetById(pairModel.Token0Erc20Id)
		})
		out2 := myUtil.MyAsync(func() erc20.ModelERC20 {
			return erc20.GetById(pairModel.Token1Erc20Id)
		})
		out3 := myUtil.MyAsync(func() dex.ModelDex {
			return dex.GetById(pairModel.DexId)
		})
		pairModel.Token0Erc20 = <-out1
		pairModel.Token1Erc20 = <-out2
		pairModel.ModelDex = <-out3
	}

	return &pairModel, nil
}

func ToString(modelPair types.ModelPair) string {
	var returnValue = fmt.Sprintf("PairId: %d ", modelPair.PairId)
	returnValue = returnValue + fmt.Sprintf("DexId: %d ", modelPair.DexId)
	returnValue = returnValue + fmt.Sprintf("PairIndex: %d ", modelPair.PairIndex)
	returnValue = returnValue + fmt.Sprintf("PairContractAddress: %s ", modelPair.PairContractAddress.String())
	returnValue = returnValue + fmt.Sprintf("Token0Erc20Id: %d ", modelPair.Token0Erc20Id)
	returnValue = returnValue + fmt.Sprintf("Token1Erc20Id: %d ", modelPair.Token1Erc20Id)
	returnValue = returnValue + fmt.Sprintf("Token0Address: %s ", modelPair.Token0Address.String())
	returnValue = returnValue + fmt.Sprintf("Token1Address: %s ", modelPair.Token1Address.String())
	returnValue = returnValue + fmt.Sprintf("Token0: %s ", erc20.ToString(modelPair.Token0Erc20))
	returnValue = returnValue + fmt.Sprintf("Token1: %s ", erc20.ToString(modelPair.Token1Erc20))
	returnValue = returnValue + fmt.Sprintf("Token0Reserves: %s ", modelPair.Token0Reserves.String())
	returnValue = returnValue + fmt.Sprintf("Token1Reserves: %s ", modelPair.Token1Reserves.String())
	returnValue = returnValue + fmt.Sprintf("ShouldFindArb: %t ", modelPair.ShouldFindArb)
	return returnValue

}

func PlsPairHasHighAmountOfPls(pair types.ModelPair, minAmountOfPls *big.Int) bool {
	
	if pair.Token0Erc20Id == wplsModelErc20.Erc20Id {
		if pair.Token0Reserves.Cmp(minAmountOfPls) >= 0 {
			return true
		}
	} else if pair.Token1Erc20Id == wplsModelErc20.Erc20Id {
		if pair.Token1Reserves.Cmp(minAmountOfPls) >= 0 {
			return true
		}
	} else {
		panic("No pls pair found for dexId: " + strconv.Itoa(pair.PairId))
	}

	return false
}
