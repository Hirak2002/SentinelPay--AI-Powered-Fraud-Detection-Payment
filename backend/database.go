package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type DatabaseConnection struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewDatabaseConnection() (*DatabaseConnection, error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DatabaseConnection{db: db}, nil
}

func (dc *DatabaseConnection) Close() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return dc.db.Close()
}

func (dc *DatabaseConnection) ExecuteQuery(query string, args ...interface{}) (sql.Result, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.db.Exec(query, args...)
}

func (dc *DatabaseConnection) QueryRow(query string, args ...interface{}) *sql.Row {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.db.QueryRow(query, args...)
}

func (dc *DatabaseConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.db.Query(query, args...)
}

func (dc *DatabaseConnection) GetDB() *sql.DB {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.db
}
