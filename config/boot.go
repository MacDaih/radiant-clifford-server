package config

import (
	"log"
	"os"
	"strings"
)

type config struct {
	port string // PORT

	serverAddr string   // MQTT server address
	clientID   string   // MQTT client add
	userName   string   // MQTT user name
	pwd        string   // MQTT user password
	topics     []string // MQTT topics

	dbName string // DB_NAME
	dbHost string // DB_HOST
	dbPort string // DB_PORT
}

var cfg config

func setValue(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("failed to boot : %s is not set", key)
	}
	return value
}

func setValues(key string) []string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("failed to boot : %s is not set", key)
	}

	return strings.Split(value, ",")
}

func Boot() {
	cfg.port = setValue("PORT")

	cfg.serverAddr = setValue("SERVER_ADDR")

	cfg.clientID = setValue("CLIENT_ID")
	cfg.userName = setValue("USER_NAME")
	cfg.pwd = setValue("PASSWORD")

	cfg.dbName = setValue("DB_NAME")
	cfg.dbHost = setValue("DB_HOST")
	cfg.dbPort = setValue("DB_PORT")
}

func GetClientID() string {
	return cfg.clientID
}

func GetUserName() string {
	return cfg.userName
}

func GetUserPasswd() string {
	return cfg.pwd
}

func GetTopics() []string {
	return cfg.topics
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
