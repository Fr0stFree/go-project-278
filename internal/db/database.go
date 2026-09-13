// Package db opens a GORM PostgreSQL connection.
package db

import (
	"database/sql"
	"fmt"

	"shortener/internal/config"
	"shortener/internal/db/models/link"
	"shortener/internal/db/models/linkvisit"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database owns the SQL connection pool and its repositories.
type Database struct {
	pool      *sql.DB
	Link      *link.Repository
	LinkVisit *linkvisit.Repository
}

// New opens PostgreSQL, and configures the connection pool.
func New(cfg *config.DataBase) (*Database, error) {
	gormDB, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	pool, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	pool.SetMaxOpenConns(cfg.MaxOpenConnections)
	pool.SetMaxIdleConns(cfg.MaxIdleConnections)
	pool.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

	return &Database{
		pool:      pool,
		Link:      link.NewRepository(gormDB),
		LinkVisit: linkvisit.NewRepository(gormDB),
	}, nil
}

// Close releases all connections in the SQL pool.
func (d *Database) Close() error {
	if err := d.pool.Close(); err != nil {
		return fmt.Errorf("close PostgreSQL connection pool: %w", err)
	}

	return nil
}
