package config

import (
	"log"
	"os"
)

type config struct {
	port string // PORT

	serverAddr string // Server Address

	dbName string // DB_NAME
	dbHost string // DB_HOST
	dbPort string // DB_PORT
}

var cfg config

func set(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("failed to boot : %s is not set", key)
	}
	return value
}

func Boot() {
	cfg.port = set("PORT")

	cfg.serverAddr = set("SERVER_ADDR")

	cfg.dbName = set("DB_NAME")
	cfg.dbHost = set("DB_HOST")
	cfg.dbPort = set("DB_PORT")
}

func GetPort() string {
	return cfg.port
}

func GetServerAddr() string {
	return cfg.serverAddr
}

func GetDBEnv() (string, string, string) {
	return cfg.dbName, cfg.dbHost, cfg.dbPort
}
