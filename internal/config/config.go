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
	storePath	string
	port            int
}

func (c *Config) GetStoreType() string {
	return c.storeType
}

func (c *Config) GetStorePath() string {
	return c.storePath 
}

func (c *Config) GetPort() int {
	return c.port
}

func (c *Config) getEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("enviorment variable is not set %s", key)
	}
	return value, nil
}

func NewConfig() *Config {
	var config Config
	portStr, err := config.getEnv("API_PORT")
	if err != nil {
		log.Fatal(err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("invalid port number provided in API_PORT")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config file: %s", viper.ConfigFileUsed())
	}

	storeType := viper.GetString("api.store.type")

	var storePath string
	switch storeType {
	case "CSV":
		storePath, err = config.getEnv("CSV_STORE_PATH")
		if err != nil {
			log.Fatal(err)
		}
	case "SQLITE":
		storePath, err = config.getEnv("SQLITE_STORE_PATH")
		if err != nil {
			log.Fatal(err)
		}
	}

	config.port = port
	config.storeType = storeType
	config.storePath = storePath

	return &config
}
