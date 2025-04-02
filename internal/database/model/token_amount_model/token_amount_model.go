package token_amount_model

import (
	"errors"
	"fmt"
	"sync"

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

func GetByOwnerAddress(ownerAddress string, limit int, offset int) ([]types.ModelTokenAmount, error) {
	db := database.GetDBConnection()

	addr := common.HexToAddress(ownerAddress)
	addrBytes := addr.Bytes()

	query := fmt.Sprintf("SELECT * FROM %s WHERE OWNER_ADDRESS = $1 ORDER BY AMOUNT DESC LIMIT $2 OFFSET $3", 
		tableName)

	rows, err := db.Query(query, addrBytes, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	results := make([]types.ModelTokenAmount, 0)
	for rows.Next() {
		tokenAmount, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, *tokenAmount)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

func GetByTokenAddress(tokenAddress common.Address, limit int, offset int) ([]types.ModelTokenAmount, error) {
	db := database.GetDBConnection()

	// Convert hex string address to bytes for PostgreSQL


	// Build query with pagination
	query := fmt.Sprintf("SELECT * FROM %s WHERE TOKEN_ADDRESS = $1 ORDER BY AMOUNT DESC LIMIT $2 OFFSET $3", 
		tableName)

	// Execute query with parameters
	rows, err := db.Query(query, tokenAddress.Bytes(), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	// Collect results
	results := make([]types.ModelTokenAmount, 0)
	for rows.Next() {
		tokenAmount, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, *tokenAmount)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	// Return empty slice if no results found
	if len(results) == 0 {
		return []types.ModelTokenAmount{}, nil
	}

	return results, nil
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
			`)

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
    
	batchSize := myConfig.GetInstance().TokenAmountModelBatchSize
	workerCount := myConfig.GetInstance().TokenAmountModelWorkerCount

	if len(tokenAddressOwnerAddressList) == 0 {
        return make(map[common.Address]map[common.Address]*types.ModelTokenAmount), nil
    }

    // Initialize result map with mutex for concurrent access
    result := make(map[common.Address]map[common.Address]*types.ModelTokenAmount)
    var resultMutex sync.RWMutex
    for _, pair := range tokenAddressOwnerAddressList {
        result[pair.TokenAddress] = make(map[common.Address]*types.ModelTokenAmount)
    }

   
    numBatches := (len(tokenAddressOwnerAddressList) + batchSize - 1) / batchSize

    // Error channel to collect errors from goroutines
    errChan := make(chan error, numBatches)
    var wg sync.WaitGroup

    // Create a channel to control the number of workers
    semaphore := make(chan struct{}, workerCount)

    // Process batches with limited concurrency
    for i := 0; i < len(tokenAddressOwnerAddressList); i += batchSize {
        wg.Add(1)
        semaphore <- struct{}{} // Acquire semaphore
        
        go func(start int) {
            defer wg.Done()
            defer func() { <-semaphore }() // Release semaphore

            end := start + batchSize
            if end > len(tokenAddressOwnerAddressList) {
                end = len(tokenAddressOwnerAddressList)
            }

            // Get new db connection for this goroutine
            db := database.GetDBConnection()

            // Build query for current batch
            var sqlBuilder strings.Builder
            sqlBuilder.WriteString("SELECT " + tokenAmountColumnNames + " FROM " + tableName + 
                " WHERE (TOKEN_ADDRESS, OWNER_ADDRESS) IN (")

            // Build the value list for current batch
            args := make([]interface{}, 0, (end-start)*2)
            for j, pair := range tokenAddressOwnerAddressList[start:end] {
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
                errChan <- fmt.Errorf("batch %d-%d error: %w", start, end, err)
                return
            }
            defer rows.Close()

            // Process results
            batchResults := make(map[common.Address]map[common.Address]*types.ModelTokenAmount)
            for rows.Next() {
                tokenAmount, err := scan(rows)
                if err != nil {
                    errChan <- fmt.Errorf("batch %d-%d scan error: %w", start, end, err)
                    return
                }

                // Initialize inner map if needed
                if _, exists := batchResults[tokenAmount.TokenAddress]; !exists {
                    batchResults[tokenAmount.TokenAddress] = make(map[common.Address]*types.ModelTokenAmount)
                }
                batchResults[tokenAmount.TokenAddress][tokenAmount.OwnerAddress] = tokenAmount
            }

            if err = rows.Err(); err != nil {
                errChan <- fmt.Errorf("batch %d-%d rows error: %w", start, end, err)
                return
            }

            // Merge batch results into main result map
            resultMutex.Lock()
            for tokenAddr, ownerMap := range batchResults {
                for ownerAddr, amount := range ownerMap {
                    result[tokenAddr][ownerAddr] = amount
                }
            }
            resultMutex.Unlock()

        }(i)
    }

    // Wait for all goroutines to finish
    go func() {
        wg.Wait()
        close(errChan)
    }()

    // Check for any errors
    for err := range errChan {
        if err != nil {
            return nil, err
        }
    }

    return result, nil
}

func scan(rows orm.Scannable) (*types.ModelTokenAmount, error) {
	tokenAmount := types.ModelTokenAmount{}

	var tempAmount []uint8

	err := rows.Scan(
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
