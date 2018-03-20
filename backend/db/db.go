package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.wdf.sap.corp/I334816/ipl18/backend/config"
)

var DB *sql.DB

func init() {
	log.Println("Initializing database...")
	dbConf := config.GetDBConfig()
	connStr := fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`, dbConf.DBUser, dbConf.DBPassword, dbConf.Host, dbConf.Port, dbConf.DBName)
	var err error

	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Println("error opening connection to db", err.Error())
		os.Exit(1)
	}
}
