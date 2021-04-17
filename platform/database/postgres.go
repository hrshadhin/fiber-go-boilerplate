package database

import (
	"fmt"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/config"
	_ "github.com/jackc/pgx/v4/stdlib" // load pgx driver for PostgreSQL
	"github.com/jmoiron/sqlx"
)

// DB holds the database
type DB struct{ *sqlx.DB }

// database instance
var defaultDB = &DB{}

// connect sets the db client of database using configuration
func (db *DB) connect(cfg *config.DB) (err error) {
	dbURI := fmt.Sprintf("host=%s port=%d sslmode=%s user=%s password=%s dbname=%s",
		cfg.Host,
		cfg.Port,
		cfg.SslMode,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

	db.DB, err = sqlx.Connect("pgx", dbURI)
	if err != nil {
		return err
	}

	// connection pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetConnMaxLifetime(cfg.MaxConnLifetime)

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return fmt.Errorf("can't sent ping to database, %w", err)
	}

	return nil
}

// GetDB returns db instance
func GetDB() *DB {
	return defaultDB
}

// ConnectDB sets the db client of database using default configuration
func ConnectDB() error {
	return defaultDB.connect(config.DBCfg())
}
