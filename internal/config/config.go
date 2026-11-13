package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	AppEnv         string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	JWTSecret      string
	JWTExpiryHours int
	AdminEmail     string
	AdminPassword  string
	AdminName      string
	ClientOrigin   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found or could not be loaded, using environment variables")
	}

	expiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "72"))
	if err != nil {
		expiryHours = 72
	}

	return &Config{
		Port:           getEnv("PORT", "8080"),
		AppEnv:         getEnv("APP_ENV", "development"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "0711"),
		DBName:         getEnv("DB_NAME", "compro_tromagnon_db"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "supersecretjwtkey_tromagnon_2026"),
		JWTExpiryHours: expiryHours,
		AdminEmail:     getEnv("ADMIN_EMAIL", "admin@tromagnon.com"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", "admin123"),
		AdminName:      getEnv("ADMIN_NAME", "Tromagnon Admin"),
		ClientOrigin:   getEnv("CLIENT_ORIGIN", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
