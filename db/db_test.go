package db_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ExtraWhy/internal-libs/db"
	"github.com/ExtraWhy/internal-libs/models/player"
)

func Test_Db_Players(t *testing.T) {

	db_sql := db.DBSqlConnection{}
	if err := db_sql.Init("sqlite3", "players.db"); err != nil {
		os.Exit(-1)
	}
	defer db_sql.Deinit()
	if err := db_sql.CreatePlayersTable(); err != nil {
		os.Exit(-1)
	}

	for i := 0; i < 10; i++ {
		test_name := fmt.Sprintf("Test-name-%d", i)
		p := player.Player[uint64]{Id: uint64(i), Money: 9999,
			CB: player.CB_Reserved{DailyLimit: 20, TotalWonDaily: 0, PointsForReward: 100, Name: test_name}}
		db_sql.AddPlayer(&p)
	}

	players := db_sql.DisplayPlayers()
	for i := 0; i < len(players); i++ {
		t.Log(players[i])
	}
}

func Test_Db_Recovery(t *testing.T) {
	db_sql := db.DBSqlConnection{}
	db_sql_pl := db.DBSqlConnection{}

	if err := db_sql_pl.Init("sqlite3", "players.db"); err != nil {
		os.Exit(-1)
	}

	if err := db_sql.Init("sqlite3", "recovery.db"); err != nil {
		os.Exit(-1)
	}
	defer db_sql.Deinit()
	defer db_sql_pl.Deinit()

	if err := db_sql_pl.CreatePlayersTable(); err != nil {
		os.Exit(-1)
	}

	for i := 0; i < 10; i++ {
		test_name := fmt.Sprintf("Test-name-%d", i)
		p := player.Player[uint64]{Id: uint64(i), Money: 9999,
			CB: player.CB_Reserved{DailyLimit: 20, TotalWonDaily: 0, PointsForReward: 100, Name: test_name}}
		db_sql_pl.AddPlayer(&p)
	}

	players := db_sql_pl.DisplayPlayers()
	if err := db_sql.CreateRecoveryTable(nil); err != nil {
		os.Exit(-1)
	}

	for i := 0; i < len(players); i++ {
		if _, err := db_sql.AddRecoveryRecord(&players[i], players[i].CB); err != nil {
			break
		}
	}
}
