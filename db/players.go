package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/ExtraWhy/internal-libs/models/player"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"go.mongodb.org/mongo-driver/bson"
)

func (db *DBSqlConnection) CreatePlayersTable() error {
	playerstable := `CREATE TABLE players (
		"id" integer ,		
		"money" integer,
		"cb_reserved" TEXT		
	  );` // SQL Statement for Create Table

	statement, err := db.db.Prepare(playerstable) // Prepare SQL Statement
	if err != nil {
		return errors.New("failed to prepare players table")
	}
	statement.Exec() // Execute SQL Statements //check for error
	return nil
}

func (db *DBSqlConnection) AddRecoveryRecord(p *player.Player[uint64], fe_state any) (int64, error) {
	feStateJSON, err := json.Marshal(fe_state)
	if err != nil {
		return -1, fmt.Errorf("failed to marshal fe_state: %w", err)
	}
	stmt, err := db.db.Prepare(`UPDATE recovery SET game_id = ?, fe_state = ? WHERE player_id = ?`)
	if err != nil {
		return -1, fmt.Errorf("failed to prepare update: %w", err)
	}

	res, err := stmt.Exec(0xff, string(feStateJSON), p.Id)
	if err != nil {
		return -1, fmt.Errorf("failed to update recovery state: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	return rowsAffected, nil
}

func (db *DBSqlConnection) CreateRecoveryTable(p *player.Player[uint64]) error {
	rectable := `CREATE TABLE recovery (
		"player_id" integer ,		
		"game_id" integer,
		"fe_state" TEXT		
	  );` // SQL Statement for Create Table

	statement, err := db.db.Prepare(rectable) // Prepare SQL Statement
	if err != nil {
		return errors.New("failed to prepare players table")
	}
	statement.Exec() // Execute SQL Statements //check for error
	return nil

}

func (db *DBSqlConnection) DisplayPlayers() []player.Player[uint64] {
	row, err := db.db.Query("SELECT * FROM players")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var list []player.Player[uint64]
	for row.Next() { // Iterate and fetch the records from result cursor
		var p player.Player[uint64]
		var cb string
		row.Scan(&p.Id, &p.Money, &cb)
		json.Unmarshal([]byte(cb), &p.CB)
		list = append(list, p)
	}
	return list
}

func (db *DBSqlConnection) AddPlayer(p *player.Player[uint64]) bool {
	if p == nil {
		return false
	}

	b, _ := json.Marshal(p.CB)
	pquery := `INSERT INTO players(id, money, cb_reserved) VALUES (?, ?, ?)`
	statement, err := db.db.Prepare(pquery)
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(p.Id, p.Money, b)
	if err != nil {
		log.Fatalln(err.Error())
		return false
	}
	return true
}

func (db *DBSqlConnection) UpdatePlayerMoney(p *player.Player[uint64]) (int64, error) {

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

func (db *DBSqlConnection) CasinoBetUpdatePlayer(p *player.Player[uint64]) (int64, error) {
	cbJSON, err := json.Marshal(p.CB)
	if err != nil {
		return -1, fmt.Errorf("failed to marshal CB: %w", err)
	}
	stmt, err := db.db.Prepare(`UPDATE players SET cb_reserved = ? WHERE id = ?`)
	if err != nil {
		return -1, fmt.Errorf("failed to prepare SQL: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(string(cbJSON), p.Id)
	if err != nil {
		return -1, fmt.Errorf("failed to update player CB: %w", err)
	}
	return res.RowsAffected()
}

func (db *NoSqlConnection) CreatePlayersTable() error {

	return nil
}

func (db *NoSqlConnection) AddPlayer(p *player.Player[uint64]) bool {

	if db.lck.TryLock() {
		defer db.lck.Unlock()
		coll := db.db.Collection("players")

		result, err := coll.InsertOne(context.TODO(), p)
		if err != nil {
			return false
		}
		return result.Acknowledged
	}
	return false
}

func (db *NoSqlConnection) DisplayPlayers() []player.Player[uint64] {
	if db.lck.TryLock() {
		defer db.lck.Unlock()
		col := db.db.Collection("players")
		cursor, err := col.Find(context.TODO(), bson.M{})
		if err != nil {
			return nil
		}
		var players []player.Player[uint64]
		for cursor.Next(context.TODO()) {
			var elem player.Player[uint64]
			err := cursor.Decode(&elem)
			if err == nil {
				players = append(players, elem)
			}
		}
		return players
	}
	return []player.Player[uint64]{} //no players
}

func (db *NoSqlConnection) UpdatePlayerMoney(p *player.Player[uint64]) (int64, error) {
	if db.lck.TryLock() {
		defer db.lck.Unlock()
		updt := bson.M{"$set": bson.M{"money": p.Money}}
		res, err := db.db.Collection("players").UpdateOne(context.TODO(), bson.M{"id": p.Id}, updt)
		if err != nil {
			return -1, errors.New("failed to update")
		}
		return res.ModifiedCount, nil
	}
	return -1, errors.New("failed to acquire lock")
}

// todo
func (db *NoSqlConnection) CasinoBetUpdatePlayer(p *player.Player[uint64]) (int64, error) {
	if db.lck.TryLock() {
		defer db.lck.Unlock()
		updt := bson.M{"$set": bson.M{"cb_reserved": p.CB}}
		res, err := db.db.Collection("players").UpdateOne(context.TODO(), bson.M{"id": p.Id}, updt)
		if err != nil {
			return -1, errors.New("failed to update player for casinobet")
		}

		return res.ModifiedCount, nil
	}
	return -1, errors.New("failed to acquire lock ")
}

func (db *NoSqlConnection) CreateRecoveryTable(p *player.Player[uint64]) error {

	if db.lck.TryLock() {
		defer db.lck.Unlock()
		updt := bson.M{"player_id": p.Id}
		_, err := db.db.Collection("recovery").InsertOne(context.TODO(), updt)

		if err != nil {
			return errors.New("failed to update recovery state ")
		}
	}
	return errors.New("failed to acquire lock")
}

func (db *NoSqlConnection) AddRecoveryRecord(p *player.Player[uint64], fe_state any) (int64, error) {
	if db.lck.TryLock() {
		defer db.lck.Unlock()
		updt := bson.M{"$set": bson.M{"game_id": 0xff, "fe_state": fe_state}}
		_, err := db.db.Collection("recovery").UpdateOne(context.TODO(), bson.M{"player_id": p.Id}, updt)
		if err != nil {
			return -1, errors.New("failed to update recovery state")
		}

		return 1, nil
	}
	return 0, errors.New("failed to acquire lock")
}
