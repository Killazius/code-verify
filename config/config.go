package config

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

func GetContainer() string {
	return os.Getenv("CONTAINER")
}

func GetEndpoint() string {
	return os.Getenv("ENDPOINT")
}
