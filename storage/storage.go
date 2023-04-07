package storage

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Storage struct {
	mongo     *mongo.Client
	dbName    string
	dbAddress string
}

func (s *Storage) GetMongo() *mongo.Client {
	return s.mongo
}

func (s *Storage) SetMongo(mongo *mongo.Client) {
	s.mongo = mongo
}

func InitStorage(dbname string, mongoAddr string) (*Storage, error) {
	var s = &Storage{
		dbName:    dbname,
		dbAddress: mongoAddr,
	}
	mongo, err := ConnectMongo(dbname, mongoAddr)
	if err != nil {
		return nil, err
	}
	s.SetMongo(mongo)
	return s, nil
}

func ConnectMongo(dbName string, mongoAddr string) (*mongo.Client, error) {
	wc := writeconcern.New(writeconcern.W(1), writeconcern.J(true))
	err := mgm.SetDefaultConfig(nil, dbName, options.Client().ApplyURI(mongoAddr).SetWriteConcern(wc))
	if err != nil {
		return nil, err
	}
	_, client, _, err := mgm.DefaultConfigs()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *Storage) CreateIndex() error {
	err := s.createIndexerStatesIndex()
	if err != nil {
		return err
	}
	err = s.CreateNameIndex()
	if err != nil {
		return err
	}
	return nil
}
