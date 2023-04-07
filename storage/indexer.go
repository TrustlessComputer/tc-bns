package storage

import (
	"bnsportal/models"
	"context"
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func (s *Storage) GetIndexerState(ctx context.Context, indexer string) (*models.IndexerState, error) {
	state := &models.IndexerState{}
	err := mgm.Coll(state).First(bson.M{
		"indexer": indexer,
	}, state)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (s *Storage) UpdateIndexerState(ctx context.Context, state *models.IndexerState) error {
	updateDoc := bson.M{
		"$set": bson.M{
			"last_indexed_block": state.LastIndexedBlock,
			"last_indexed_at":    state.LastIndexedAt,
			"state":              state.State,
		},
	}

	_, err := mgm.Coll(state).UpdateOne(ctx, bson.M{"indexer": state.Indexer}, updateDoc)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateIndexerState(ctx context.Context, state *models.IndexerState) error {
	err := mgm.Coll(state).Create(state)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) createIndexerStatesIndex() error {
	startTime := time.Now()
	model := []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{Key: "indexer", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := mgm.Coll(&models.IndexerState{}).Indexes().CreateMany(context.Background(), model)
	if err != nil {
		log.Printf("failed to index coins in %v", time.Since(startTime))
		return err
	}

	return nil
}
