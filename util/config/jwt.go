package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func JwtKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("JWT_SECRET_KEY")
}