package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	JWTSecretKey   string
	AdminSecretKey string
	UserSecretKey  string
)

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: No .env file found, using defaults")
	}

	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	AdminSecretKey = os.Getenv("ADMIN_SECRET_KEY")
	UserSecretKey = os.Getenv("USER_SECRET_KEY")

	if JWTSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY is required in .env file")
	}
	if AdminSecretKey == "" {
		log.Fatal("ADMIN_SECRET_KEY is required in .env file")
	}
	if UserSecretKey == "" {
		log.Fatal("USER_SECRET_KEY is required in .env file")
	}
}
