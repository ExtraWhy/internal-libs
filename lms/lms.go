package lms

import "github.com/ExtraWhy/internal-libs/models/player"

// limit management system for casinobet only
// shall write to db, ours or theirs its a TBD
// otherwise we can hold in mongo atlas test account those for testing purposes
// see player updated struct
type LMS struct {
	players []player.Player[uint64]
}

// TODO: make function see test
/*
 */

func (l *LMS) find_by_id(id uint64) *player.Player[uint64] {
	for i := 0; i < len(l.players); i++ {
		if l.players[i].Id == id {
			return &l.players[i]
		}
	}
	return nil
}

func (l *LMS) UpdatePlayerLimitById(id, lim uint64) bool {
	if p := l.find_by_id(id); p != nil {
		p.CB.DailyLimit = lim
		return true
	}
	return false
}

func (l *LMS) UpdateDailyReward(id, rew uint64) bool {
	if p := l.find_by_id(id); p != nil {
		p.CB.PointsForReward = rew
		return true
	}
	return false
}
