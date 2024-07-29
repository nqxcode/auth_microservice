package config

import (
	"github.com/joho/godotenv"
)

// Load env parameters from file by path
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
