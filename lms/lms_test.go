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
		l.players = append(l.players, player.Player{Id: v + 1})
	}
	return &l, nil
}

func PrinMe(l *LMS) {
	for i := 0; i < len(l.players); i++ {
		fmt.Printf("Player id [%d] -> stats [%d][%d] \r\n", l.players[i].Id, l.players[i].DailyLimit, l.players[i].PointsForReward)
	}
}

func TestLMS(t *testing.T) {
	lms, _ := MakeLmsTest()
	PrinMe(lms)
	lms.UpdateDailyReward(2, 10)
	lms.UpdatePlayerLimitById(6, 10)
	dbiface := &db.NoSqlConnection{}
	//test account used here
	dbiface.Init("Cluster0", "cryptowincryptowin:EfK0weUUe7t99Djx")
	for i := 1; i < 10; i++ {
		dbiface.CasinoBetUpdatePlayer(&lms.players[i-1])
	}
	PrinMe(lms)
}
