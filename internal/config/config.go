package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
)

type Config struct {
	storeType       string
	storePathCSV    string
	storePathSQLite string
	port            int
}

func (c *Config) GetStoreType() string {
	return c.storeType
}

func (c *Config) GetStorePathCSV() string {
	return c.storePathCSV
}

func (c *Config) GetStorePathSqlite() string {
	return c.storePathSQLite
}

func (c *Config) GetPort() int {
	return c.port
}

func getEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("enviorment variable is not set %s", key)
	}
	return value, nil
}

func NewConfig() *Config {
	var config Config
	portStr, err := getEnv("API_PORT")
	if err != nil {
		log.Fatal(err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("invalid port number provided in API_PORT")
	}

	storePathCSV, err := getEnv("CSV_STORE_PATH")
	if err != nil {
		log.Fatal(err)
	}
	storePathSqlite, err := getEnv("SQLITE_STORE_PATH")
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config file: %s", viper.ConfigFileUsed())
	}

	storeType := viper.GetString("api.store.type")

	config.port = port
	config.storeType = storeType
	config.storePathCSV = storePathCSV
	config.storePathSQLite = storePathSqlite

	return &config
}
