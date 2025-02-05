package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JWTSecretKey string

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: No .env file found, using defaults")
	}

	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	if JWTSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY is required in .env file")
	}
	log.Println("JWT Secret Key Loaded: ", JWTSecretKey)
}

