package db

import (
	"log"

	"github.com/ExtraWhy/internal-libs/models/player"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func (db *DBConnection) CreatePlayersTable() {
	playerstable := `CREATE TABLE players (
		"id" integer ,		
		"money" integer,
		"name" TEXT		
	  );` // SQL Statement for Create Table

	statement, err := db.db.Prepare(playerstable) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
}

func (db *DBConnection) DisplayPlayers() []player.Player {
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

func (db *DBConnection) AddPlayer(p *player.Player) bool {
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
