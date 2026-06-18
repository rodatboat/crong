package database

import (
	"log"

	"github.com/rodatboat/crong/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDb initializes the database connection and runs auto-migrations
func InitDb(cfg *config.Config) *gorm.DB {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Note: Schema migrations are handled by Goose (./migrations)
	// GORM is used only for querying, not schema management
	DB = db
	return db
}
