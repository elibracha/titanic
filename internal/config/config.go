package config

import (
	"fmt"
	"os"
	"strconv"
	"titanic-api/internal/passenger"
)

type Config struct {
	storeType       passenger.StoreType
	storePathCSV    string
	storePathSQLite string
	port            int
}

func (c *Config) GetStoreType() passenger.StoreType {
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

func NewConfig() (*Config, error) {
	var config Config

	storeType, err := getEnv("STORE_TYPE")
	if err != nil {
		return nil, err
	}

	portStr, err := getEnv("API_PORT")
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port number provided in API_PORT")
	}

	storePathCSV, err := getEnv("CSV_STORE_PATH")
	if err != nil {
		return nil, err
	}
	storePathSqlite, err := getEnv("SQLITE_STORE_PATH")
	if err != nil {
		return nil, err
	}

	config.port = port
	config.storeType = passenger.StoreType(storeType)
	config.storePathCSV = storePathCSV
	config.storePathSQLite = storePathSqlite

	return &config, nil
}
