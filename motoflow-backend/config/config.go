package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBDriver    string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	DBPath      string // for sqlite
	JWTSecret   string
	Environment string
}

var AppConfig *Config

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite" // default to sqlite for instant out-of-the-box local execution; postgres enabled via env
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "motoflow-ahass-secret-key-nusantara-sakti-2026"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "motoflow_ahass"
	}

	dbSSL := os.Getenv("DB_SSLMODE")
	if dbSSL == "" {
		dbSSL = "disable"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "motoflow.db"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	AppConfig = &Config{
		Port:        port,
		DBDriver:    driver,
		DBHost:      dbHost,
		DBPort:      dbPort,
		DBUser:      dbUser,
		DBPassword:  dbPass,
		DBName:      dbName,
		DBSSLMode:   dbSSL,
		DBPath:      dbPath,
		JWTSecret:   jwtSecret,
		Environment: env,
	}

	return AppConfig
}

