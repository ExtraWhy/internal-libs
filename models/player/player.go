package player

import "github.com/google/uuid"

type SpecializedID interface {
	uint64 | uuid.UUID
}

type CB_Barfill struct {
	GameBars [6]int `json:"bars"`
}

// player structure for playing a game
type Player[T SpecializedID] struct {
	Id    T      `json:"id"`
	Money uint64 `json:"money"`
	//related for main/casinobet/dev
	DailyLimit      uint64     `json:"daily_limit"`
	TotalWonDaily   uint64     `json:"total_won_daily"`
	PointsForReward uint64     `json:"points_for_reward"`
	Name            string     `json:"name"`
	BarFill         CB_Barfill `json:"bar_fill"`
}
