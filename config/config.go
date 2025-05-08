package config

import (
	"fmt"
	"log"
	"os"
)

// Provide easy way to maintain adding/removing env variables
const (
	SYMBOL_ENV  = "SYMBOL_ENV"
	NDAYS_ENV   = "NDAYS_ENV"
	APIKEY_PATH = "APIKEY_PATH"
)

// Read apikey file path from env variable and extract apikey from file
func ReadAPIKey() (string, error) {
	apikeyPath := os.Getenv(APIKEY_PATH)
	log.Printf("Reading apikey file: %s", apikeyPath)

	// Check pathname
	_, err := os.Stat(apikeyPath)
	if err != nil {
		log.Printf("Error locating apikey path: %s", err.Error())
		return "", fmt.Errorf("Unable to locate apikey path")
	}

	apikey, err := os.ReadFile(apikeyPath)
	if err != nil {
		log.Printf("Error reading apikey: %s", err.Error())
		return "", fmt.Errorf("Unable to get valid apikey")
	}

	return string(apikey), nil
}

// Return provided env variable value, if it exists
func GetEnvData(envStr string) (string, error) {
	log.Printf("Looking up ENV: %s", envStr)
	data, exist := os.LookupEnv(envStr)
	if !exist {
		log.Printf("%s ENV is Unset", envStr)
		return "", fmt.Errorf("Unable to locate ENV entry: %s", envStr)
	}
	return data, nil
}
