package lms

import (
	"fmt"
	"testing"

	"github.com/ExtraWhy/internal-libs/db"
	"github.com/ExtraWhy/internal-libs/models/player"
)

func MakeLmsTest() (*LMS, error) {
	l := LMS{}
	for i := 0; i < 10; i++ {
		var v uint64 = (uint64)(i)
		l.players = append(l.players, player.Player[uint64]{Id: v + 1})
	}
	return &l, nil
}

func PrintMe(l *LMS) {
	for i := 0; i < len(l.players); i++ {
		fmt.Printf("Player id [%d] -> stats [%d][%d] \r\n", l.players[i].Id, l.players[i].DailyLimit, l.players[i].PointsForReward)
	}
}

func TestLMS(t *testing.T) {
	lms, _ := MakeLmsTest()
	PrintMe(lms)
	lms.UpdateDailyReward(2, 1)
	lms.UpdatePlayerLimitById(6, 10)
	lms.UpdatePlayerLimitById(8, 100)
	lms.UpdatePlayerLimitById(10, 1000)

	dbiface := &db.NoSqlConnection{}
	//test account used here
	dbiface.Init("Cluster0", "cryptowincryptowin:EfK0weUUe7t99Djx")
	defer dbiface.Deinit()
	for i := 1; i < 10; i++ {
		dbiface.CasinoBetUpdatePlayer(&lms.players[i-1])
	}
	PrintMe(lms)
}
