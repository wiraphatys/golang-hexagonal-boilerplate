package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDatabase(host, user, password, dbname, port, sslmode, TimeZone string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host,
		user,
		password,
		dbname,
		port,
		sslmode,
		TimeZone,
	)

	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // Slow SQL threshold
			LogLevel:      logger.Error, // Log level
			Colorful:      true,         // Disable color
		},
	)

	gormConfig := &gorm.Config{
		Logger:                                   logger,
		DisableForeignKeyConstraintWhenMigrating: false,
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database during init: %v. DSN (sensitive parts might be shown): %s", err, dsn)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("FATAL: Failed to get underlying sql.DB from GORM during init: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("FATAL: Failed to ping database during init: %v", err)
	}

	log.Println("INFO: Database connection established successfully during init.")

	return gormDB
}
