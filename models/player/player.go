package player

import (
	"github.com/google/uuid"
)

type SpecializedID interface {
	uint64 | uuid.UUID
}

type CB_Reserved struct {
	DailyLimit      uint64 `json:"daily_limit" bson:"daily_limit"`
	TotalWonDaily   uint64 `json:"total_won_daily" bson:"total_won_daily"`
	PointsForReward uint64 `json:"points_for_reward" bson:"points_for_reward"`
	Name            string `json:"name" bson:"name"`
	BarFill         [6]int `json:"bar_fill" bson:"bar_fill"`
}

type Player[T SpecializedID] struct {
	Id    T           `json:"id" bson:"id"`
	Money uint64      `json:"money" bson:"money"`
	CB    CB_Reserved `json:"cb_reserved" bson:"cb_reserved"`
}

func (p *Player[T]) PlayerCBSchema() string {
	ret := "{\"daily_limit\": 0,\"total_won_daily\":0,\"points_for_reward\":0,\"name\":\"uname\",\"bar_fill\":[0,0,0,0,0,0]}"
	return ret
}
