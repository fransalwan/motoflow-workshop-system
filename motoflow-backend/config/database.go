package config

import (
	"fmt"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"motoflow-backend/models"
)

var DB *gorm.DB

func InitDB(cfg *Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if cfg.DBDriver == "postgres" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
		
		log.Printf("[Database] Connecting to PostgreSQL on %s:%s/%s...", cfg.DBHost, cfg.DBPort, cfg.DBName)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err != nil {
			log.Printf("[Database WARNING] PostgreSQL connection failed: %v. Falling back to SQLite (%s) for local offline operation.", err, cfg.DBPath)
			db, err = gorm.Open(sqlite.Open(cfg.DBPath), gormConfig)
		}
	} else {
		log.Printf("[Database] Using SQLite storage: %s", cfg.DBPath)
		db, err = gorm.Open(sqlite.Open(cfg.DBPath), gormConfig)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	log.Println("[Database] Running auto-migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Category{},
		&models.Sparepart{},
		&models.Mechanic{},
		&models.WorkOrder{},
		&models.WorkOrderItem{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	DB = db
	log.Println("[Database] Migrations applied successfully.")
	return DB, nil
}

