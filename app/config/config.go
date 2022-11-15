package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	ServerPort       string
	DropboxAuthToken string
	PostgresPort     string
	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgesDBName    string
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	return &Config{
		ServerPort:       os.Getenv("SERVER_PORT"),
		DropboxAuthToken: os.Getenv("DROPBOX_AUTH_TOKEN"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgesDBName:    os.Getenv("POSTGRES_DB_NAME"),
	}
}
