package db

import (
	"context"
	"errors"
	"log"

	"github.com/ExtraWhy/internal-libs/models/player"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"go.mongodb.org/mongo-driver/bson"
)

func (db *DBSqlConnection) CreatePlayersTable() error {
	playerstable := `CREATE TABLE players (
		"id" integer ,		
		"money" integer,
		"name" TEXT		
	  );` // SQL Statement for Create Table

	statement, err := db.db.Prepare(playerstable) // Prepare SQL Statement
	if err != nil {
		return errors.New("failed to prepare players table")
	}
	statement.Exec() // Execute SQL Statements //check for error
	return nil
}

func (db *DBSqlConnection) DisplayPlayers() []player.Player {
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

func (db *NoSqlConnection) CreatePlayersTable() error {

	return nil
}

func (db *DBSqlConnection) AddPlayer(p *player.Player) bool {
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

func (db *DBSqlConnection) UpdatePlayerMoney(p *player.Player) (int64, error) {

	if p == nil {
		return -1, errors.New("nil reference to player")
	}
	udquery := `UPDATE players SET money = ? WHERE id = ?;`
	if result, err := db.db.Exec(udquery, p.Money, p.Id); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (db *NoSqlConnection) AddPlayer(p *player.Player) bool {

	coll := db.db.Collection("players")

	result, err := coll.InsertOne(context.TODO(), p)
	if err != nil {
		return false
	}
	return result.Acknowledged
}

func (db *NoSqlConnection) DisplayPlayers() []player.Player {
	col := db.db.Collection("players")
	cursor, err := col.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil
	}
	var players = make([]player.Player, 5)
	for cursor.Next(context.TODO()) {
		var elem player.Player
		err := cursor.Decode(&elem)
		if err == nil {
			players = append(players, elem)
		}
	}
	return players
}

func (db *NoSqlConnection) UpdatePlayerMoney(p *player.Player) (int64, error) {
	updt := bson.M{"$set": bson.M{"money": p.Money}}
	res, err := db.db.Collection("players").UpdateOne(context.TODO(), bson.M{"id": p.Id}, updt)
	if err != nil {
		return -1, errors.New("failed to updated")
	}
	return res.ModifiedCount, nil
}
