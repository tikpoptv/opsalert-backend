package db

import (
	"database/sql"
	"fmt"
	"opsalert/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes database connection
func InitDB() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.DBConfig.Host,
		config.AppConfig.DBConfig.Port,
		config.AppConfig.DBConfig.User,
		config.AppConfig.DBConfig.Password,
		config.AppConfig.DBConfig.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	return nil
}

// CloseDB closes database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
