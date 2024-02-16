package utils

import (
	"log"
	"os"
)

func SetEnv(key, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
