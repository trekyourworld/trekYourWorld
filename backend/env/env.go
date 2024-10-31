package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	envFile := os.Getenv("GODOTENV")
	if envFile == "" {
		envFile = ".env"
	}

	if err := godotenv.Load(envFile); err != nil {
		return fmt.Errorf("error loading %s file", envFile)
	}
	return nil
}
