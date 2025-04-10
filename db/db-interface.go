package db

import "fmt"

type DbIface interface {
	Init(driver string, dsn string) error
	Deinit() error
}

type UnimplementedDbConnector struct {
}

func (UnimplementedDbConnector) mustEmbedUnimplementedDbConnector() {}

func (UnimplementedDbConnector) Init(driver string, dsn string) error {
	return fmt.Errorf("Must implement method Init")
}

func (UnimplementedDbConnector) Deinit() error {
	return fmt.Errorf("Must implement method Deinit")
}

func (UnimplementedDbConnector) Connect() error {
	return fmt.Errorf("Must implement method Deinit")
}
