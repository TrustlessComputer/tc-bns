package models

import "github.com/kamva/mgm/v3"

type RegisteredNameInfo struct {
	mgm.DefaultModel  `bson:",inline"`
	Owner             string `json:"owner" bson:"owner"`
	ID                string `json:"id" bson:"id"`
	Name              string `json:"name" bson:"name"`
	RegisteredAtBlock uint64 `json:"registered_at_block" bson:"registered_at_block"`
}
