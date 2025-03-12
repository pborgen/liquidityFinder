package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"sync"

	"github.com/pborgen/liquidityFinder/internal/myConfig"
)

var lock = &sync.Mutex{}
var dbConnection *sql.DB
var cfg = myConfig.GetInstance()


func GetDBConnection() *sql.DB {

	if dbConnection == nil {
		lock.Lock()
		defer lock.Unlock()

		var host string
		var port int
		var user string
		var password string
		var dbName string 
		var sslMode string
		if cfg.UseLocalDB {
			log.Info().Msg("Using local db")
			host = "localhost"
			port = 5432
			user = "postgres"
			password = "<your-password>"
			dbName = "<your-db-name>"
			sslMode = "disable"
		} else {
			log.Info().Msg("Using remote db")
			host = cfg.PostgresHost
			port = cfg.PostgresPort
			user = cfg.PostgresUser
			password = cfg.PostgresPassword
			dbName = cfg.PostgresDB
			sslMode = cfg.PostgresSSLMode
		}

		if dbConnection == nil {

			pSqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
				"password=%s dbname=%s sslmode=%s binary_parameters=yes",
				host, port, user, password, dbName, sslMode)
			
			log.Info().Msg(pSqlInfo)

			db, err := sql.Open("postgres", pSqlInfo)

			if err != nil {
				panic("Could not get db connection" + err.Error())
			}

			db.SetMaxIdleConns(3)
			db.SetConnMaxLifetime(3 * time.Minute)
			db.SetMaxOpenConns(30)
			db.SetConnMaxIdleTime(1 * time.Minute)

			dbConnection = db

			err = db.Ping()
			if err != nil {
				panic("Could not ping db" + err.Error())
			}
		}
	}

	return dbConnection

}
