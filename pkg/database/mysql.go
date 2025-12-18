package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Printf("Warning: .env.%s file not found, using system environment variables", env)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	_ = godotenv.Load(".env." + env)

	if os.Getenv("DB_USER") == "" || os.Getenv("DB_NAME") == "" || os.Getenv("DB_HOST") == "" || os.Getenv("DB_PORT") == "" || os.Getenv("DB_PASS") == "" {
		log.Fatal("Database environment variables are not set")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	var gormConfig = &gorm.Config{}
	if env == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan instance database: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database tidak merespon (Ping gagal): %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Koneksi database berhasil dan terverifikasi!")
	return db
}
