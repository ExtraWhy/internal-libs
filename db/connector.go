package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/ExtraWhy/internal-libs/player"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type DBconnection struct {
	//more to be filled here and implement interfaces for other db

	db *sql.DB // priv connection
}

func (db *DBconnection) Init(dbname string) bool {
	os.Remove(dbname)

	file, err := os.Create(dbname) // Create SQLite file todo - as argument
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()

	db.db, _ = sql.Open("sqlite3", dbname) // Open the created SQLite File
	createTable(db.db)

	//todo checks
	return true
}

func (db *DBconnection) Deinit() {
	db.db.Close()
}

func (db *DBconnection) DisplayPlayers() []player.Player {
	row, err := db.db.Query("SELECT * FROM players")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var list []player.Player
	for row.Next() { // Iterate and fetch the records from result cursor
		var p player.Player
		row.Scan(&p.Id, &p.Money, &p.Name)
		list = append(list, p)
	}
	return list
}

func (db *DBconnection) AddPlayer(p *player.Player) bool {
	if p == nil {
		return false
	}
	pquery := `INSERT INTO players(id, money, name) VALUES (?, ?, ?)`
	statement, err := db.db.Prepare(pquery)
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(p.Id, p.Money, p.Name)
	if err != nil {
		log.Fatalln(err.Error())
		return false
	}
	return true
}

func createTable(db *sql.DB) {
	playerstable := `CREATE TABLE players (
		"id" integer ,		
		"money" integer,
		"name" TEXT		
	  );` // SQL Statement for Create Table

	statement, err := db.Prepare(playerstable) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
}
