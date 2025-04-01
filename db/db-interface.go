package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type TableInitializer func(db *sql.DB) error

type DBConnection struct {
	//more to be filled here and implement interfaces for other db

	db *sql.DB // priv connection
}

func (dbc *DBConnection) Init(dbname string, initializers []TableInitializer) error {
	os.Remove(dbname)

	file, err := os.Create(dbname) // Create SQLite file todo - as argument
	if err != nil {
		log.Fatal(err.Error())
		return fmt.Errorf("failed to create db file: %w", err)
	}
	file.Close()

	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	dbc.db = db

	for _, initializer := range initializers {
		if err := initializer(dbc.db); err != nil {
			return fmt.Errorf("Failed to create table: %w", err)
		}
	}

	return nil
}

func (db *DBConnection) Deinit() {
	db.db.Close()
}
