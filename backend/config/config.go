package config

import (
	"encoding/json"
	"log"
	"os"

	"github.wdf.sap.corp/I334816/ipl18/backend/models"
)

var serverConfig models.Config

var GetConfig = func() models.Config {
	return serverConfig
}

var GetDBConfig = func() models.DBConfig {
	return serverConfig.DB
}

var GetHashConfig = func() models.HashConfig {
	return serverConfig.Hash
}

var GetJwtConfig = func() models.JWTConfig {
	return serverConfig.JWT
}

func init() {
	log.Println("Parsing config...")
	f, err := os.Open("./config.json")
	if err != nil {
		log.Println("error opening config file", err.Error())
		os.Exit(1)
	}

	err = json.NewDecoder(f).Decode(&serverConfig)
	if err != nil {
		log.Println("error parsing config file", err.Error())
		os.Exit(1)
	}
}
