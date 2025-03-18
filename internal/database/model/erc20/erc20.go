package erc20

import (
	"database/sql"
	"fmt"
	"unicode"

	"strconv"

	"time"
	"unicode/utf8"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/database"
	"github.com/pborgen/liquidityFinder/internal/database/model/orm"
	"github.com/pborgen/liquidityFinder/myConst"
	"github.com/rs/zerolog/log"
)

type ModelERC20 struct {
	Erc20Id         int            `postgres.Table:"ERC20_ID" json:"erc20_id"`
	NetworkId       int            `postgres.Table:"NETWORK_ID" json:"network_id"`
	ContractAddress common.Address `postgres.Table:"CONTRACT_ADDRESS" json:"contract_address"`
	Name            string         `postgres.Table:"NAME" json:"name"`
	Symbol          string         `postgres.Table:"SYMBOL" json:"symbol"`
	Decimal         uint8          `postgres.Table:"DECIMAL" json:"decimal"`
	ShouldFindArb   bool           `postgres.Table:"SHOULD_FIND_ARB" json:"should_find_arb"`
	IsValidated     bool           `postgres.Table:"IS_VALIDATED" json:"is_validated"`
	IsTaxToken      bool           `postgres.Table:"IS_TAX_TOKEN" json:"is_tax_token"`
	TaxPercentage   float64        `postgres.Table:"TAX_PERCENTAGE" json:"tax_percentage"`
	ProcessedIsTaxToken bool       `postgres.Table:"PROCESSED_IS_TAX_TOKEN" json:"processed_is_tax_token"`
	SymbolImageUrl  string `postgres.Table:"SYMBOL_IMAGE_URL" json:"symbol_image_url"`
}

const primaryKey = "ERC20_ID"
const tableName = "ERC20"

var columnNames = orm.GetColumnNames(ModelERC20{})



func Insert(erc20Model ModelERC20) (ModelERC20, error) {
	db := database.GetDBConnection()

	sqlStatement := orm.CreateInsertStatement(ModelERC20{}, tableName, primaryKey)

	id := 0
	var returnValue ModelERC20

	name := erc20Model.Name
	name = sanitizeUTF8(name)

	lengthName := utf8.RuneCountInString(name)
	if lengthName > 150 {
		name = trimToLengthRunes(name, 150)
	}

	symbol := erc20Model.Symbol

	lengthSymbol := utf8.RuneCountInString(symbol)
	if lengthSymbol > 150 {
		symbol = trimToLengthRunes(symbol, 150)
	}

	err := db.QueryRow(
		sqlStatement, 
		erc20Model.NetworkId, 
		erc20Model.ContractAddress.String(), 
		name, 
		symbol, 
		erc20Model.Decimal, 
		erc20Model.ShouldFindArb, 
		erc20Model.IsValidated,
		erc20Model.IsTaxToken,
		erc20Model.TaxPercentage,
		erc20Model.ProcessedIsTaxToken,
		erc20Model.SymbolImageUrl,
	).Scan(&id)
	if err == nil {
		returnValue = GetById(id)
	}

	return returnValue, err
}

func UpdateSymbolImageUrl(contractAddress common.Address, symbolImageUrl string) (error) {
	db := database.GetDBConnection()
	sqlStatement := `
	UPDATE ERC20 
	SET SYMBOL_IMAGE_URL=$2
	WHERE CONTRACT_ADDRESS=$1
	`
	_, err := db.Exec(sqlStatement, contractAddress.String(), symbolImageUrl)
	return err
}

func GetRandom() []ModelERC20 {
	db := database.GetDBConnection()
	sql := "SELECT " + columnNames + " FROM " + tableName + " ORDER BY RANDOM() LIMIT 10"
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	results := make([]ModelERC20, 0)
	for rows.Next() {
		erc20, err := scan(rows)
		if err != nil {
			panic(err)
		}
		results = append(results, *erc20)
	}

	return results
}

func GetAllNonProcessedIsTaxToken() []ModelERC20 {
	db := database.GetDBConnection()
	sql := "SELECT " + columnNames + " FROM " + tableName + " WHERE PROCESSED_IS_TAX_TOKEN = false"
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	results := make([]ModelERC20, 0)
	for rows.Next() {
		erc20, err := scan(rows)
		if err != nil {
			panic(err)
		}
		results = append(results, *erc20)
	}
	return results
}


func UpdateTax(contractAddress common.Address, isTaxToken bool, taxPercentage float64) (error) {
	db := database.GetDBConnection()
	sqlStatement := `
	UPDATE ERC20 
	SET IS_TAX_TOKEN=$2, TAX_PERCENTAGE=$3, PROCESSED_IS_TAX_TOKEN=true
	WHERE CONTRACT_ADDRESS=$1
	`
	_, err := db.Exec(sqlStatement, contractAddress.String(), isTaxToken, taxPercentage)

	return err
}

func GetById(id int) ModelERC20 {

	db := database.GetDBConnection()

	sql := "SELECT " + columnNames + " FROM " + tableName + " WHERE ERC20_ID = $1"
	row := db.QueryRow(sql, id)

	erc20Model, err := scan(row)

	if err != nil {
		log.Fatal().Msgf(err.Error())
		panic("Erc20 does not exists with id:" + strconv.Itoa(id))
	}

	return *erc20Model
}

func GetWPLS() (ModelERC20, error) {
	return GetByContractAddress(common.HexToAddress(myConst.WPLS_ADDRESS_STRING), 1)
}

func GetDecimalsByContractAddress(contractAddress common.Address) (int8, error) {
	db := database.GetDBConnection()

	row := db.QueryRow(
		"SELECT DECIMAL FROM "+tableName+" WHERE CONTRACT_ADDRESS = $1", contractAddress.String())

	var decimal int8
	err := row.Scan(&decimal)
	if err != nil {
		return 0, err
	}

	return decimal, nil
}

func AsyncGetDecimalsByContractAddress(contractAddress common.Address) <-chan int8 {
    resultChan := make(chan int8)
    
    go func() {
        defer close(resultChan)
		decimals, err := GetDecimalsByContractAddress(contractAddress)
		if err != nil {
			resultChan <- -1
			log.Error().Err(err).Msg("Failed to get decimals for token with address: " + contractAddress.String() + " error: " + err.Error())
			return
		}
        resultChan <- decimals
    }()
    
    return resultChan
}

func GetByContractAddress(contractAddress common.Address, retries int) (ModelERC20, error) {

	var erc20 *ModelERC20
	var err error

	for i := 0; i < retries; i++ {
		db := database.GetDBConnection()

		row := db.QueryRow(
		"SELECT "+columnNames+" FROM "+tableName+" WHERE CONTRACT_ADDRESS = $1", contractAddress.String())

		erc20, err = scan(row)
		
		if err == nil {
			break
		}
		
		log.Error().Err(err).Msg("Error in Erc20GetByContractAddress:  contractAddress: " + contractAddress.String())
		time.Sleep(1 * time.Second)
	}

	return *erc20, nil
}

func GetAll() []ModelERC20 {

	db := database.GetDBConnection()

	results := make([]ModelERC20, 0)
	rows, err := db.Query("SELECT " + columnNames + " FROM " + tableName)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	//var erc20 ERC20Model
	for rows.Next() {
		erc20, err := scan(rows)
		if err != nil {

			panic(err)
		}
		results = append(results, *erc20)
	}

	return results

}

func ExistsByContractAddress(contractAddress common.Address) bool {
	return orm.MyExists(tableName, "CONTRACT_ADDRESS", contractAddress.String())
}

func scan(rows orm.Scannable) (*ModelERC20, error) {
	erc20 := ModelERC20{}

	var contractAddressTemp string
	var symbolImageUrlTemp sql.NullString
	err := rows.Scan(
		&erc20.Erc20Id, 
		&erc20.NetworkId, 
		&contractAddressTemp, 
		&erc20.Name, 
		&erc20.Symbol, 
		&erc20.Decimal, 
		&erc20.ShouldFindArb, 
		&erc20.IsValidated, 
		&erc20.IsTaxToken, 
		&erc20.TaxPercentage, 
		&erc20.ProcessedIsTaxToken, 
		&symbolImageUrlTemp)
	if err != nil {
		return &ModelERC20{}, err
	}

	erc20.ContractAddress = common.HexToAddress(contractAddressTemp)
	if symbolImageUrlTemp.Valid {
		erc20.SymbolImageUrl = symbolImageUrlTemp.String
	} else {
		erc20.SymbolImageUrl = ""
	}
	return &erc20, nil
}

func ToString(modelErc20 ModelERC20) string {
	var returnValue = fmt.Sprintf("Erc20Id: %d ", modelErc20.Erc20Id)
	returnValue = returnValue + fmt.Sprintf("NetworkId: %d ", modelErc20.NetworkId)
	returnValue = returnValue + fmt.Sprintf("ContractAddress: %s ", modelErc20.ContractAddress.String())
	returnValue = returnValue + fmt.Sprintf("Name: %s ", modelErc20.Name)
	returnValue = returnValue + fmt.Sprintf("Symbol: %s ", modelErc20.Symbol)
	returnValue = returnValue + fmt.Sprintf("SymbolImageUrl: %s ", modelErc20.SymbolImageUrl)
	returnValue = returnValue + fmt.Sprintf("Decimal: %d ", modelErc20.Decimal)
	returnValue = returnValue + fmt.Sprintf("ShouldFindArb: %t ", modelErc20.ShouldFindArb)
	returnValue = returnValue + fmt.Sprintf("IsValidated: %t ", modelErc20.IsValidated)

	return returnValue
}

func trimRunes(s string) string {
    // Convert to runes
    runes := []rune(s)

    // Find the start (skip leading whitespace)
    start := 0
    for start < len(runes) && unicode.IsSpace(runes[start]) {
        start++
    }

    // If all whitespace, return empty string
    if start >= len(runes) {
        return ""
    }

    // Find the end (skip trailing whitespace)
    end := len(runes) - 1
    for end >= start && unicode.IsSpace(runes[end]) {
        end--
    }

    // Return the trimmed substring
    return string(runes[start : end+1])
}

func sanitizeUTF8(s string) string {
    if utf8.ValidString(s) {
        return s
    }
    return string([]rune(s))
}

func trimToLengthRunes(s string, maxLength int) string {
    // If already short enough, return as-is
    if utf8.RuneCountInString(s) <= maxLength {
        return s
    }

    // Convert to runes and trim
    runes := []rune(s)
    return string(runes[:maxLength])
}
