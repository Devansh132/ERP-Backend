package database

import (
	"fmt"
	"log"

	"school-erp-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error
	var dialector gorm.Dialector

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
	)

	switch config.AppConfig.DBDriver {
	case "postgres":
		dialector = postgres.Open(dsn)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.AppConfig.DBUser,
			config.AppConfig.DBPassword,
			config.AppConfig.DBHost,
			config.AppConfig.DBPort,
			config.AppConfig.DBName,
		)
		dialector = mysql.Open(dsn)
	default:
		return fmt.Errorf("unsupported database driver: %s", config.AppConfig.DBDriver)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func AutoMigrate() error {
	// Import models here to avoid circular dependency
	// This will be called from main.go after importing models
	return nil
}



