package db

import (
	"fmt"

	"github.com/ExtraWhy/internal-libs/models/player"
)

type DbIface[T player.SpecializedID] interface {
	Init(driver string, dsn string) error
	Deinit() error
	UpdatePlayerMoney(p *player.Player[T]) (int64, error)
	DisplayPlayers() []player.Player[T]
	AddPlayer(p *player.Player[T]) bool
	CreatePlayersTable() error
	CasinoBetUpdatePlayer(*player.Player[T]) (int64, error) //only for casino bet client

	CreateRecoveryTable(p *player.Player[T]) error
	AddRecoveryRecord(*player.Player[T], any) (int64, error)
}

type UnimplementedDbConnector struct {
}

func (UnimplementedDbConnector) mustEmbedUnimplementedDbConnector() {}

func (UnimplementedDbConnector) CreatePlayersTable() error {
	return fmt.Errorf("Must implement method Init")
}

func (UnimplementedDbConnector) Init(driver string, dsn string) error {
	return fmt.Errorf("Must implement method Init")
}

func (UnimplementedDbConnector) Deinit() error {
	return fmt.Errorf("Must implement method Deinit")
}
