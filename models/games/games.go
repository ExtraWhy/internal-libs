package games

// game pair will hold all relevant data needed for the games
// for recovery and user state like how much games are
// also we can have the rtp as an info and a hash code for the random number (TBD)
// reserved for casinobet and future uses
type GameData struct {
	ID          uint32 `json:"id"`
	FreeGames   uint32 `json:"freegames"`
	TotalPlayed uint64 `json:"num_played"`
	TotalWins   uint64 `json:"total_wins"`
	Reserved    any    `json:"reserved"`
}

type Games struct {
	LastPlayedGame uint32     `json:"last_played"`
	GameList       []GameData `json:"game_list"`
}
