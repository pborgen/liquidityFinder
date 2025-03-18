package token_amount_model

import (
	"errors"

	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pborgen/liquidityFinder/internal/database"
	"github.com/pborgen/liquidityFinder/internal/database/model/orm"
	"github.com/pborgen/liquidityFinder/internal/myConfig"
	"github.com/pborgen/liquidityFinder/internal/myUtil"
	"github.com/rs/zerolog/log"

	"github.com/pborgen/liquidityFinder/internal/types"
)

// Table schema:
// CREATE TABLE TOKEN_AMOUNT (
//     ID SERIAL PRIMARY KEY,                                    -- Auto-incrementing integer
//     TOKEN_ADDRESS BYTEA NOT NULL,                            -- Ethereum addresses stored as 20 bytes
//     OWNER_ADDRESS BYTEA NOT NULL,                            -- Ethereum addresses stored as 20 bytes
//     AMOUNT NUMERIC(78,0) NOT NULL,                           -- For large token amounts
//     LAST_BLOCK_NUMBER_UPDATED BIGINT NOT NULL,               -- Range: -9,223,372,036,854,775,808 to +9,223,372,036,854,775,807
//     LAST_LOG_INDEX_UPDATED INTEGER NOT NULL,                 -- Range: -2,147,483,648 to +2,147,483,647
//     CONSTRAINT token_amount_unique UNIQUE (TOKEN_ADDRESS, OWNER_ADDRESS)
// );

const primaryKey = "ID"
const tableName = "TOKEN_AMOUNT"

var tokenAmountModelInsertBatchSize = myConfig.GetInstance().TokenAmountModelInsertBatchSize
var tokenAmountColumnNames = orm.GetColumnNames(types.ModelTokenAmount{})

func init() {

}

func BatchInsertOrUpdate(tokenAmounts []types.ModelTokenAmount) (error) {
	if len(tokenAmounts) == 0 {
		return nil
	}


	for i := 0; i < len(tokenAmounts); i += tokenAmountModelInsertBatchSize {
		end := i + tokenAmountModelInsertBatchSize
		if end > len(tokenAmounts) {
			end = len(tokenAmounts)
		}

		db := database.GetDBConnection()

		var sqlBuilder strings.Builder
		sqlBuilder.WriteString(`
			INSERT INTO ` + tableName + ` (
				TOKEN_ADDRESS, OWNER_ADDRESS, AMOUNT, LAST_BLOCK_NUMBER_UPDATED, LAST_LOG_INDEX_UPDATED
			) VALUES `)

		args := []interface{}{}
		count := 1

		for _, element := range tokenAmounts[i:end] {
			if !myUtil.IsWithinNumeric78Range(element.Amount) {
				continue
			}

			if count > 1 {
				sqlBuilder.WriteString(",")
			}

			sqlBuilder.WriteString("($")
			sqlBuilder.WriteString(strconv.Itoa(count))
			sqlBuilder.WriteString(", $")
			sqlBuilder.WriteString(strconv.Itoa(count+1))
			sqlBuilder.WriteString(", $")
			sqlBuilder.WriteString(strconv.Itoa(count+2))
			sqlBuilder.WriteString(", $")
			sqlBuilder.WriteString(strconv.Itoa(count+3))
			sqlBuilder.WriteString(", $")
			sqlBuilder.WriteString(strconv.Itoa(count+4))
			sqlBuilder.WriteString(")")

			count += 5

			// Store raw bytes of addresses
			args = append(
				args, 
				element.TokenAddress,
				element.OwnerAddress,
				element.Amount.String(),
				element.LastBlockNumberUpdated,
				element.LastLogIndexUpdated,
			)
		}

		sqlBuilder.WriteString(`
			ON CONFLICT (TOKEN_ADDRESS, OWNER_ADDRESS) DO UPDATE SET 
				AMOUNT = EXCLUDED.AMOUNT,
				LAST_BLOCK_NUMBER_UPDATED = EXCLUDED.LAST_BLOCK_NUMBER_UPDATED,
				LAST_LOG_INDEX_UPDATED = EXCLUDED.LAST_LOG_INDEX_UPDATED
			RETURNING ID`)

		sql := sqlBuilder.String()

		// Prepare the insert statement
		_, err := db.Exec(sql, args...)
		
		if err != nil {
			log.Error().Err(err).Msg("Error executing batch insert")
			panic(err)
		}

	}

	return nil
}

func GetLargestLastBlockNumberUpdated() (uint64, error) {
	db := database.GetDBConnection()

	count, err := GetCount()
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 0, nil
	}

	rows := db.QueryRow("SELECT MAX(LAST_BLOCK_NUMBER_UPDATED) FROM " + tableName)

	if rows == nil {
		return 0, errors.New("could not find smallest block number")
	}

	var smallestBlockNumber uint64
	err = rows.Scan(&smallestBlockNumber)
	if err != nil {
		return 0, err
	}

	return smallestBlockNumber, nil
}

func GetCount() (int, error) {
	db := database.GetDBConnection()

	rows := db.QueryRow("SELECT COUNT(*) FROM " + tableName)

	var count int
	err := rows.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetById(id int) (*types.ModelTokenAmount, error) {

	db := database.GetDBConnection()

	rows := db.QueryRow(
		"SELECT "+tokenAmountColumnNames+" FROM "+tableName+" WHERE ID = $1", id)

	if rows == nil {
		return nil, errors.New("Could not find pairId:" + strconv.Itoa(id))
	}

	tokenAmount, err := scan(rows)

	if err != nil {
		return nil, err
	}

	return tokenAmount, nil
}

// Remove empty addresses from slice
func removeEmptyAddresses(addresses []common.Address) []common.Address {
	result := make([]common.Address, 0, len(addresses))
	emptyAddress := common.Address{}
	
	for _, addr := range addresses {
		// Skip zero/empty addresses
		if addr != emptyAddress {
			result = append(result, addr)
		}
	}
	return result
}

func GetByContractAddressAndOwner(contractAddressList []common.Address, ownerList []common.Address) (map[common.Address]map[common.Address]*types.ModelTokenAmount, error) {
	// Clean input arrays by removing empty addresses
	contractAddressList = removeEmptyAddresses(contractAddressList)
	ownerList = removeEmptyAddresses(ownerList)

	if len(contractAddressList) == 0 || len(ownerList) == 0 {
		return make(map[common.Address]map[common.Address]*types.ModelTokenAmount), nil
	}

	db := database.GetDBConnection()

	// Build query with IN clause for both lists
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("SELECT " + tokenAmountColumnNames + " FROM " + tableName + " WHERE TOKEN_ADDRESS = ANY($1) AND OWNER_ADDRESS = ANY($2)")

	// Convert addresses to byte arrays for postgres
	contractAddressBytes := make([][]byte, len(contractAddressList))
	ownerAddressBytes := make([][]byte, len(ownerList))

	for i, addr := range contractAddressList {
		contractAddressBytes[i] = addr.Bytes()
	}

	for i, addr := range ownerList {
		ownerAddressBytes[i] = addr.Bytes()
	}

	sql := sqlBuilder.String()

	// Execute query with byte arrays
	rows, err := db.Query(sql, pq.Array(contractAddressBytes), pq.Array(ownerAddressBytes))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize empty map
	result := make(map[common.Address]map[common.Address]*types.ModelTokenAmount)

	// Pre-initialize inner maps for all contract addresses
	for _, contractAddr := range contractAddressList {
		result[contractAddr] = make(map[common.Address]*types.ModelTokenAmount)
	}

	// Iterate through rows and build the map
	for rows.Next() {
		tokenAmount, err := scan(rows)
		if err != nil {
			return nil, err
		}
		
		result[tokenAmount.TokenAddress][tokenAmount.OwnerAddress] = tokenAmount
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func scan(rows orm.Scannable) (*types.ModelTokenAmount, error) {
	tokenAmount := types.ModelTokenAmount{}

	var tempAmount []uint8

	err := rows.Scan(
		&tokenAmount.Id,
		&tokenAmount.TokenAddress,
		&tokenAmount.OwnerAddress,
		&tempAmount,
		&tokenAmount.LastBlockNumberUpdated,
		&tokenAmount.LastLogIndexUpdated,
	)

	if err != nil {
		return &types.ModelTokenAmount{}, err
	}

	tempAmountString := string(tempAmount)

	tempAmountBigInt, success := new(big.Int).SetString(tempAmountString, 10)
	if !success {
		return &types.ModelTokenAmount{}, errors.New("failed to convert string to big.Int")
	}

	tokenAmount.Amount = tempAmountBigInt

	return &tokenAmount, nil
}
