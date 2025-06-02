package recovery

import (
	"fmt"

	"github.com/ExtraWhy/internal-libs/models/games"
)

type MapUserGames map[string]games.Games

var sGames = MapUserGames{}

func AddRecovery(usr string, gameid, free uint32, won uint64) {
	if v, ok := sGames[usr]; ok {
		for i := 0; i < len(v.GameList); i++ {
			if v.GameList[i].ID == gameid {
				v.LastPlayedGame = gameid
				if free > 0 {
					v.GameList[i].FreeGames += free
				}
				if won > 0 {
					v.GameList[i].TotalWins++
				}
				v.GameList[i].TotalPlayed++
			}
		}
	}
}

// test only for now !!!
func DumpRecovery() {
	fmt.Println(sGames)
}
