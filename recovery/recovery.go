package recovery

import (
	"github.com/ExtraWhy/internal-libs/models/games"
)

var sGames = games.Games{}

//ID          uint32 `json:"id"`
//	FreeGames   uint32 `json:"freegames"`
//	TotalPlayed uint64 `json:"num_played"`
//	TotalWins   uint64 `json:"total_wins"`

func AddToRecovery(gameid, freegames uint32, totalplayed, wins uint64) {
	gd := games.GameData{ID: gameid, FreeGames: freegames, TotalPlayed: totalplayed, TotalWins: wins}
	sGames.GameList = append(sGames.GameList, gd)
	sGames.LastPlayedGame = gameid
}

func GetGameById(gameid uint32) *games.GameData {
	for i := 0; i < len(sGames.GameList); i++ {
		if sGames.GameList[i].ID == gameid {
			return &sGames.GameList[i]
		}
	}
	return nil
}
