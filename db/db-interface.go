package db

import (
	"fmt"

	"github.com/ExtraWhy/internal-libs/models/player"
)

type DbIface interface {
	Init(driver string, dsn string) error
	Deinit() error
	UpdatePlayerMoney(p *player.Player) (int64, error)
	DisplayPlayers() []player.Player
	AddPlayer(p *player.Player) bool
	CreatePlayersTable() error
	CasinoBetUpdatePlayer(*player.Player) error //only for casino bet client
}

type UnimplementedDbConnector struct {
}

func (UnimplementedDbConnector) mustEmbedUnimplementedDbConnector() {}

func (UnimplementedDbConnector) CreatePlayersTable() error {
	return fmt.Errorf("Must implement method Init")
}

func (UnimplementedDbConnector) AddPlayer(p *player.Player) bool {
	return false
}

func (UnimplementedDbConnector) UpdatePlayerMoney(p *player.Player) (int64, error) {
	return -1, fmt.Errorf("Must implement method Init")
}

// used only for casinobet client
func (UnimplementedDbConnector) CasinoBetUpdatePlayer(*player.Player) (int64, error) {
	return -1, fmt.Errorf("Must implement method CasinoBetUpdatePlayer")
}

func (UnimplementedDbConnector) DisplayPlayers() []player.Player {
	return nil
}

func (UnimplementedDbConnector) Init(driver string, dsn string) error {
	return fmt.Errorf("Must implement method Init")
}

func (UnimplementedDbConnector) Deinit() error {
	return fmt.Errorf("Must implement method Deinit")
}
