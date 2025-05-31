package player

// player structure for playing a game
type Player struct {
	Id    uint64 `json:"id"`
	Money uint64 `json:"money"`
	Name  string `json:"name"`
}
