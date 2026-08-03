// Package postgres opens a GORM PostgreSQL connection.
package postgres

import (
	"fmt"
	"shortener/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DataBase wraps the configured GORM database connection.
type DataBase struct {
	config *config.DataBase
	DB     *gorm.DB
}

// NewDataBase opens PostgreSQL, runs AutoMigrate for the provided models, and configures the connection pool.
func NewDataBase(cfg *config.DataBase, models ...any) (*DataBase, error) {
	sslMode := "disable"
	if cfg.IsSSLEnabled {
		sslMode = "require"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	dbSQL, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	dbSQL.SetMaxOpenConns(cfg.MaxOpenConnections)
	dbSQL.SetMaxIdleConns(cfg.MaxIdleConnections)
	dbSQL.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

	return &DataBase{
		config: cfg,
		DB:     db,
	}, nil
}
