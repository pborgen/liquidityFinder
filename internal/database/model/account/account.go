package account

import (
	"github.com/pborgen/liquidityFinder/internal/database"
	"github.com/pborgen/liquidityFinder/internal/database/model/orm"
	"github.com/pborgen/liquidityFinder/internal/myEncrypt"
)

type ModelAccount struct {
	AccountId   int    `postgres.Table:"ACCOUNT_ID"`
	Name        string `postgres.Table:"NAME"`
	Description string `postgres.Table:"DESCRIPTION"`
	PublicKey   string `postgres.Table:"PUBLIC_KEY"`
	PrivateKey  string `postgres.Table:"PRIVATE_KEY"`
	IsActive    bool   `postgres.Table:"IS_ACTIVE"`
}

const primaryKey = "ACCOUNT_ID"
const tableName = "ACCOUNT"


func Insert(name string, description string, publicKey string, privateKeyInTheClear string, isActive bool) (*ModelAccount, error) {
	db := database.GetDBConnection()

	privateKey := encryptPrivateKey(privateKeyInTheClear)	

	sqlStatement := orm.CreateInsertStatement(ModelAccount{}, tableName, primaryKey)

	id := 0
	var returnValue *ModelAccount

	err := db.QueryRow(
		sqlStatement, 
		name, 
		description, 
		publicKey, 
		privateKey, 
		isActive,
	).Scan(&id)
	if err == nil {
		returnValue, err = GetById(id)
		if err != nil {
			return nil, err
		}
	}

	return returnValue, err
}

func GetById(id int) (*ModelAccount, error) {

	var account ModelAccount

	db := database.GetDBConnection()

	err := db.QueryRow(
		"SELECT * FROM ACCOUNT WHERE account_id = $1", id).Scan(
			&account.AccountId, 
			&account.Name, 
			&account.Description, 
			&account.PublicKey, 
			&account.PrivateKey, 
			&account.IsActive,
		)

	if err != nil {
		return nil, err
	}

	decryptPrivateKey(&account)


	return &account, nil
}

func GetAllActive() []ModelAccount {
	all := GetAll()

	active := make([]ModelAccount, 0)

	for _, account := range all {
		if account.IsActive {
			active = append(active, account)
		}
	}

	return active
}

func GetAll() []ModelAccount {

	db := database.GetDBConnection()

	results := make([]ModelAccount, 0)
	rows, err := db.Query("SELECT account_id, name, description, public_key, private_key, is_active FROM ACCOUNT")

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var account ModelAccount
	for rows.Next() {
		err := rows.Scan(&account.AccountId, &account.Name, &account.Description, &account.PublicKey, &account.PrivateKey, &account.IsActive)
		if err != nil {
			panic(err)
		}

		decryptPrivateKey(&account)
		
		results = append(results, account)
	}

	return results

}

func decryptPrivateKey(account *ModelAccount) {
	privateKey, err := myEncrypt.Decrypt(account.PrivateKey)
	if err != nil {
		panic(err)
	}
	account.PrivateKey = privateKey
}

func encryptPrivateKey(privateKeyInTheClear string) string {
	privateKey, err := myEncrypt.Encrypt(privateKeyInTheClear)
	if err != nil {
		panic(err)
	}
	return privateKey
}