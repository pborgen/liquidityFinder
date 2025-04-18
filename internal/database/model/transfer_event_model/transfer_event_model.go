package transfer_event_model

import (
	"errors"

	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/database"
	"github.com/pborgen/liquidityFinder/internal/database/model/orm"
	"github.com/pborgen/liquidityFinder/internal/myUtil"

	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/rs/zerolog/log"
)


const tableName = "TRANSFER_EVENT"


var transferEventColumnNames = orm.GetColumnNames(types.ModelTransferEvent{})



func init() {

}

func GetAllForAddressGroupBy(address common.Address, viewMode string) ([]types.TransferEventGroupBy, error) {

	db := database.GetDBConnection()

	
	selectStatement := "SELECT FROM_ADDRESS, TO_ADDRESS, COUNT(*) as TRANSACTION_COUNT "
	fromAddressStatement := " FROM " + tableName
	whereStatement := ""
	groupByStatement := " GROUP BY FROM_ADDRESS, TO_ADDRESS "
	orderByStatement := " ORDER BY TRANSACTION_COUNT DESC"
	limitStatement := " LIMIT 5000"

	if viewMode == "all" {
		whereStatement = " WHERE FROM_ADDRESS = $1 OR TO_ADDRESS = $1"
	} else if viewMode == "in" {
		whereStatement = " WHERE TO_ADDRESS = $1"
	} else if viewMode == "out" {
		whereStatement = " WHERE FROM_ADDRESS = $1"
	}

	sql := selectStatement + fromAddressStatement + whereStatement + groupByStatement + orderByStatement + limitStatement

	rows, err := db.Query(
		sql, 
		address.Bytes())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []types.TransferEventGroupBy
	for rows.Next() {
		var event types.TransferEventGroupBy
		err := rows.Scan(&event.FromAddress, &event.ToAddress, &event.TransactionCount)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil	
}
	
func GetAllForAddress(address common.Address, limit int, offset int) ([]types.ModelTransferEvent, error) {
	db := database.GetDBConnection()

	rows, err := db.Query(
		"SELECT " + transferEventColumnNames + " FROM " + tableName + " WHERE FROM_ADDRESS = $1 OR TO_ADDRESS = $1 LIMIT $2 OFFSET $3", 
		address.Hex(), 
		limit, 
		offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []types.ModelTransferEvent
	for rows.Next() {
		event, err := scan(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}	

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}	

func Exists(blockNumber uint64, index int) (bool, error) {
	db := database.GetDBConnection()

	rows := db.QueryRow("SELECT COUNT(1) FROM " + tableName + " WHERE BLOCK_NUMBER = $1 AND INDEX = $2", blockNumber, index)

	var count int
	err := rows.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetByBlockNumberAndIndex(blockNumber uint64, index int) (*types.ModelTransferEvent, error) {
	db := database.GetDBConnection()

	rows := db.QueryRow("SELECT " + transferEventColumnNames + " FROM " + tableName + " WHERE BLOCK_NUMBER = $1 AND INDEX = $2", blockNumber, index)
	
	if rows == nil {
		return nil, errors.New("transfer event not found")
	}

	transferEvent, err := scan(rows)
	if err != nil {
		return nil, err
	}

	return transferEvent, nil
}

func GetEventsForBlockRangeOrdered(fromBlock uint64, toBlock uint64) ([]types.ModelTransferEvent, error) {
	db := database.GetDBConnection()

	rows, err := db.Query("SELECT " + transferEventColumnNames + " FROM " + tableName + " WHERE BLOCK_NUMBER >= $1 AND BLOCK_NUMBER <= $2 ORDER BY BLOCK_NUMBER ASC, LOG_INDEX ASC", fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []types.ModelTransferEvent
	for rows.Next() {
		event, err := scan(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func GetEventsForBlockRangeAsMap(fromBlock uint64, toBlock uint64) (map[struct{ BlockNumber uint64; LogIndex int }]*types.ModelTransferEvent, error) {
	db := database.GetDBConnection()

	rows, err := db.Query("SELECT " + transferEventColumnNames + " FROM " + tableName + " WHERE BLOCK_NUMBER >= $1 AND BLOCK_NUMBER <= $2 ORDER BY BLOCK_NUMBER ASC, LOG_INDEX ASC", fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make(map[struct{ BlockNumber uint64; LogIndex int }]*types.ModelTransferEvent)
	for rows.Next() {
		event, err := scan(rows)
		if err != nil {
			return nil, err
		}
		key := struct{ BlockNumber uint64; LogIndex int }{
			BlockNumber: event.BlockNumber,
			LogIndex:    event.LogIndex,
		}
		events[key] = event
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, errors.New("no transfer events found")
	}

	return events, nil
}

func GetEventsForBlockRangeAsIndexedMap(fromBlock uint64, toBlock uint64) (map[uint64][]*types.ModelTransferEvent, error) {
	db := database.GetDBConnection()

	rows, err := db.Query("SELECT " + transferEventColumnNames + " FROM " + tableName + " WHERE BLOCK_NUMBER >= $1 AND BLOCK_NUMBER <= $2 ORDER BY BLOCK_NUMBER ASC, LOG_INDEX ASC", fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize map with empty slices for each block
	events := make(map[uint64][]*types.ModelTransferEvent)
	for rows.Next() {
		event, err := scan(rows)
		if err != nil {
			return nil, err
		}
		events[event.BlockNumber] = append(events[event.BlockNumber], event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, errors.New("no transfer events found")
	}

	return events, nil
}

func BatchInsertOrUpdate(transferEvents []types.ModelTransferEvent) ([]int, error) {
	if len(transferEvents) == 0 {
		return nil, nil
	}

	const batchSize = 4000

	for i := 0; i < len(transferEvents); i += batchSize {
		end := i + batchSize
		if end > len(transferEvents) {
			end = len(transferEvents)
		}

		db := database.GetDBConnection()

		var sqlBuilder strings.Builder
		sqlBuilder.WriteString(`
			INSERT INTO ` + tableName + ` (
				TRANSACTION_HASH, BLOCK_NUMBER, LOG_INDEX, CONTRACT_ADDRESS, 
				FROM_ADDRESS, TO_ADDRESS, EVENT_VALUE
			) VALUES `)

		args := []interface{}{}
		count := 1

		for _, event := range transferEvents[i:end] {
			if !myUtil.IsWithinNumeric78Range(event.EventValue) {
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
			sqlBuilder.WriteString(", $")
			sqlBuilder.WriteString(strconv.Itoa(count+5))
			sqlBuilder.WriteString(", $")
			sqlBuilder.WriteString(strconv.Itoa(count+6))
			sqlBuilder.WriteString(")")

			count += 7

			args = append(
				args, 
				event.TransactionHash,
				event.BlockNumber, 
				event.LogIndex, 
				event.ContractAddress, 
				event.FromAddress, 
				event.ToAddress, 
				event.EventValue.String(),
			)
		}

		sqlBuilder.WriteString(`
			ON CONFLICT (BLOCK_NUMBER, LOG_INDEX) DO UPDATE SET 
				CONTRACT_ADDRESS = EXCLUDED.CONTRACT_ADDRESS,
				FROM_ADDRESS = EXCLUDED.FROM_ADDRESS,
				TO_ADDRESS = EXCLUDED.TO_ADDRESS,
				EVENT_VALUE = EXCLUDED.EVENT_VALUE 
			`)

		sql := sqlBuilder.String()

		// Prepare the insert statement
		_, err := db.Exec(sql, args...)
		if err != nil {
			log.Error().Msg(sql)
			return nil, err
		}
	}
	return nil, nil
}

func DoAnyRowsExists() (bool, error) {
	db := database.GetDBConnection()

	rows := db.QueryRow("SELECT EXISTS (SELECT 1 FROM " + tableName + " LIMIT 1)")
	
	var exists bool
	err := rows.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func GetLargestBlockNumber() (uint64, error) {
	db := database.GetDBConnection()

	rows := db.QueryRow("SELECT MAX(BLOCK_NUMBER) FROM " + tableName)

	if rows == nil {
		return 0, errors.New("could not find largest block number")
	}

	var largestBlockNumber uint64
	err := rows.Scan(&largestBlockNumber)
	if err != nil {
		return 0, err
	}

	return largestBlockNumber, nil
}

func GetSmallestBlockNumber() (uint64, error) {
	db := database.GetDBConnection()

	count, err := GetCount()
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 0, nil
	}

	rows := db.QueryRow("SELECT MIN(BLOCK_NUMBER) FROM " + tableName)

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

func GetById(id int) (*types.ModelTransferEvent, error) {

	db := database.GetDBConnection()

	rows := db.QueryRow(
		"SELECT "+transferEventColumnNames+" FROM "+tableName+" WHERE ID = $1", id)

	if rows == nil {
		return nil, errors.New("Could not find pairId:" + strconv.Itoa(id))
	}

	transferEvent, err := scan(rows)

	if err != nil {
		return nil, err
	}

	return transferEvent, nil
}


func scan(rows orm.Scannable) (*types.ModelTransferEvent, error) {
	transferEvent := types.ModelTransferEvent{}

	var tempValue []uint8

	err := rows.Scan(
		&transferEvent.TransactionHash,
		&transferEvent.BlockNumber,
		&transferEvent.LogIndex,
		&transferEvent.ContractAddress,
		&transferEvent.FromAddress,
		&transferEvent.ToAddress,
		&tempValue,
	)

	if err != nil {
		return &types.ModelTransferEvent{}, err
	}

	tempValueString := string(tempValue)

	tempValueBigInt, success := new(big.Int).SetString(tempValueString, 10)
	if !success {
		return &types.ModelTransferEvent{}, errors.New("failed to convert string to big.Int")
	}

	transferEvent.EventValue = tempValueBigInt

	return &transferEvent, nil
}
