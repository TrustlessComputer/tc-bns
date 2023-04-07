package bnsindexer

import (
	"bnsportal/storage"
	"testing"
)

func TestIndexer_Start(t *testing.T) {
	type fields struct {
		DefaultRPC   string
		DefaultRPCWS string
		BNSContract  string
		Storage      *storage.Storage
	}
	storageInst, _ := storage.InitStorage("bfs-gateway", "mongodb://0.0.0.0:27017")
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "sdfds",
			fields: fields{
				DefaultRPC:  "https://tc-regtest.trustless.computer",
				BNSContract: "0xbA7EED9D832D0194824e00364C96C027E1Ffd221",
				Storage:     storageInst,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexer := &Indexer{
				DefaultRPC:   tt.fields.DefaultRPC,
				DefaultRPCWS: tt.fields.DefaultRPCWS,
				BNSContract:  tt.fields.BNSContract,
				Storage:      tt.fields.Storage,
			}
			indexer.Start()
		})
	}
}
