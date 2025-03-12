package namevalue

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/pborgen/liquidityFinder/internal/database"
	"github.com/pborgen/liquidityFinder/internal/database/model/orm"
	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/rs/zerolog/log"
)

var dexColumnNames = orm.GetColumnNames(types.NameValue{})

const tableName = "NAME_VALUE"

const primaryKey = "ID"

func UpdateValue(name string, value string) (error) {
	db := database.GetDBConnection()
	sqlStatement := `
	UPDATE NAME_VALUE 
	SET VALUE=$2
	WHERE NAME=$1
	`
	_, err := db.Exec(sqlStatement, name, value)

	return err
}

func Insert(nameValue types.NameValue) (types.NameValue, error) {
	db := database.GetDBConnection()

	sqlStatement := orm.CreateInsertStatement(types.NameValue{}, tableName, primaryKey)

	id := 0
	var returnValue types.NameValue

	err := db.QueryRow(
		sqlStatement, 
		nameValue.Name, 
		nameValue.Value, 
		nameValue.DataType, 
	).Scan(&id)
	if err == nil {
		returnValue, err = GetById(id)
		if err != nil {
			return types.NameValue{}, err
		}
	}

	return returnValue, err
}

func GetByName(name string) (types.NameValue, error) {

	db := database.GetDBConnection()
	query := "SELECT " + dexColumnNames + " FROM " + tableName + " WHERE NAME = $1"
	rows, err := db.Query(query, name)

	if err != nil {
		panic("Error getting NameValue with name:" + name)
	}

	defer rows.Close()

	if rows.Next() {
		nameValue, err := scan(rows)

		if err != nil {
			return types.NameValue{}, err
		}
		return *nameValue, nil
	} else {
		return types.NameValue{}, errors.New("Could not find name value with name:" + name)
	}
}

func GetById(id int) (types.NameValue, error) {

	db := database.GetDBConnection()
	query := "SELECT " + dexColumnNames + " FROM " + tableName + " WHERE ID = $1"
	rows, err := db.Query(query, id)

	if err != nil {
		return types.NameValue{}, errors.New("Error getting NameValue with id:" + strconv.Itoa(id))
	}

	defer rows.Close()

	if rows.Next() {
		nameValue, err := scan(rows)

		if err != nil {
			return types.NameValue{}, err
		}
		return *nameValue, nil
	} else {
		return types.NameValue{}, errors.New("Could not find name value with id:" + strconv.Itoa(id))
	}
}

func Exists(name string) bool {
	return orm.MyExists(tableName, "NAME", name)
}


func scan(rows *sql.Rows) (*types.NameValue, error) {
	
	nameValue := types.NameValue{}

	err := rows.Scan(
		&nameValue.Id, 
		&nameValue.Name, 
		&nameValue.Value, 
		&nameValue.DataType, 
	)
	if err != nil {
		log.Error().Msgf("Could not scan rows for name value", err)

	}

	return &nameValue, nil
}
