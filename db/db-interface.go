package db

import (
	"database/sql"
	"fmt"
	"os"
)

type TableInitializer func(db *sql.DB) error

type DBConnection struct {
	//more to be filled here and implement interfaces for other db

	db *sql.DB // priv connection
}

func (dbc *DBConnection) Init(driver string, dsn string) error {
	if err := prepareDriver(driver, dsn); err != nil {
		return err
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	dbc.db = db
	return nil
}

func (db *DBConnection) Deinit() {
	db.db.Close()
}

func prepareDriver(driver, dsn string) error {
	switch driver {
	case "sqlite3":
		_ = os.Remove(dsn)
		file, err := os.Create(dsn)
		if err != nil {
			return fmt.Errorf("failed to create db file: %w", err)
		}
		return file.Close()

	case "postgres", "mysql":
		// TO DO if needed, but no file prepr is needed so maybe skip even

	default:
		return fmt.Errorf("unsupported driver: %s", driver)
	}

	return nil
}

func (dbc *DBConnection) SetupSchema(initializers ...TableInitializer) error {
	for _, initializer := range initializers {
		if err := initializer(dbc.db); err != nil {
			return fmt.Errorf("Failed to create table: %w", err)
		}
	}
	return nil
}
