package initializers

import (
	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
