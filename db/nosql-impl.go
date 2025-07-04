package db

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type NoSqlConnection struct {
	dbc *mongo.Client
	db  *mongo.Database
	lck sync.Mutex
	UnimplementedDbConnector
}

func (dbc *NoSqlConnection) Init(driver string, dsn string) error {
	//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	//todo - move that string
	//cryptowincryptowin:EfK0weUUe7t99Djx
	//Cluster0
	if dbc.lck.TryLock() {
		defer dbc.lck.Unlock()
		login := fmt.Sprintf("mongodb+srv://%s@cluster0.w07rcmn.mongodb.net/?appName=%s", dsn, driver)
		opts := options.Client().ApplyURI(login).SetServerAPIOptions(serverAPI)
		var err error
		// Create a new client and connect to the server
		dbc.dbc, err = mongo.Connect(opts)
		if err != nil {
			return err
		}

		// Send a ping to confirm a successful connection
		if err := dbc.dbc.Ping(context.TODO(), readpref.Primary()); err != nil {
		}

		dbc.db = dbc.dbc.Database("cryptowin") // .Collection("players")

		fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	}
	return nil

}

func (dbc *NoSqlConnection) Deinit() error {
	if dbc.lck.TryLock() {
		defer dbc.lck.Unlock()
		if err := dbc.dbc.Disconnect(context.TODO()); err != nil {
			return err
		}
	}
	return nil
}
