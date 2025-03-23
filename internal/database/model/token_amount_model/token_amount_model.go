package token_amount_model

import (
	"errors"
	"fmt"

	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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

type TokenAddressOwnerAddress struct {
	TokenAddress common.Address
	OwnerAddress common.Address
}

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

	log.Info().Msgf("Inserted or Updated %d token amounts", len(tokenAmounts))

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

func GetByContractAddressAndOwner(tokenAddressOwnerAddressList []TokenAddressOwnerAddress) (map[common.Address]map[common.Address]*types.ModelTokenAmount, error) {
	if len(tokenAddressOwnerAddressList) == 0 {
		return make(map[common.Address]map[common.Address]*types.ModelTokenAmount), nil
	}

	// Initialize result map
	result := make(map[common.Address]map[common.Address]*types.ModelTokenAmount)
	for _, pair := range tokenAddressOwnerAddressList {
		result[pair.TokenAddress] = make(map[common.Address]*types.ModelTokenAmount)
	}

	// Process in batches of 1000 pairs
	batchSize := 1000
	db := database.GetDBConnection()

	for i := 0; i < len(tokenAddressOwnerAddressList); i += batchSize {
		end := i + batchSize
		if end > len(tokenAddressOwnerAddressList) {
			end = len(tokenAddressOwnerAddressList)
		}

		// Build query for current batch
		var sqlBuilder strings.Builder
		sqlBuilder.WriteString("SELECT " + tokenAmountColumnNames + " FROM " + tableName + 
			" WHERE (TOKEN_ADDRESS, OWNER_ADDRESS) IN (")

		// Build the value list for current batch
		args := make([]interface{}, 0, (end-i)*2)
		for j, pair := range tokenAddressOwnerAddressList[i:end] {
			if j > 0 {
				sqlBuilder.WriteString(",")
			}
			sqlBuilder.WriteString(fmt.Sprintf("($%d,$%d)", j*2+1, j*2+2))
			args = append(args, pair.TokenAddress.Bytes(), pair.OwnerAddress.Bytes())
		}
		sqlBuilder.WriteString(")")

		// Execute query for current batch
		rows, err := db.Query(sqlBuilder.String(), args...)
		if err != nil {
			return nil, err
		}

		// Process results
		for rows.Next() {
			tokenAmount, err := scan(rows)
			if err != nil {
				rows.Close()
				return nil, err
			}
			result[tokenAmount.TokenAddress][tokenAmount.OwnerAddress] = tokenAmount
		}
		rows.Close()

		if err = rows.Err(); err != nil {
			return nil, err
		}
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
