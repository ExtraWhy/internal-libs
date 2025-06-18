package player

// player structure for playing a game
type Player struct {
	Id    uint64 `json:"id"`
	Money uint64 `json:"money"`
	//related for main/casinobet/dev
	DailyLimit      uint64 `json:"daily_limit"`
	TotalWonDaily   uint64 `json:"total_won_daily"`
	PointsForReward uint64 `json:"points_for_reward"`
	Name            string `json:"name"`
}
