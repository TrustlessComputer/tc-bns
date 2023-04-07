package models

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type IndexerState struct {
	mgm.DefaultModel `bson:",inline"`
	Indexer          string    `json:"indexer" bson:"indexer"`
	LastIndexedBlock uint64    `json:"last_indexed_block" bson:"last_indexed_block"`
	LastIndexedAt    time.Time `json:"last_indexed_at" bson:"last_indexed_at"`
	Cron             string    `json:"cron" bson:"cron"`
	State            string    `json:"state" bson:"state"`
}
