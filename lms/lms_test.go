package lms

import (
	"fmt"
	"testing"

	"github.com/ExtraWhy/internal-libs/models/player"
)

func MakeLmsTest() (*LMS, error) {
	l := LMS{}
	for i := 0; i < 10; i++ {
		var v uint64 = (uint64)(i)
		l.players = append(l.players, player.Player{Id: v})
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
	PrinMe(lms)
}
