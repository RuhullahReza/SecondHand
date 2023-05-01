package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func DatabaseHost() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("DATABASE_HOST")
}

func DatabasePort() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("DATABASE_PORT")
}

func DatabaseUsername() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("DATABASE_USERNAME")
}

func DatabasePassword() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("DATABASE_PASSWORD")
}

func DatabaseName() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("DATABASE_NAME")
}

func DatabaseMode() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("DATABASE_SSL_MODE")
}