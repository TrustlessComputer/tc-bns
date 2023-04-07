package main

import (
	_ "bnsportal/abi"
	"bnsportal/constants"
	"bnsportal/gateway"
	"bnsportal/storage"
	bnsindexer "bnsportal/workers/bns_indexer"
	"log"
	"os"
)

func main() {

	defaultRPC := os.Getenv("DEFAULT_RPC")
	defaultRPCWS := os.Getenv("DEFAULT_RPCWS")
	port := os.Getenv("PORT")
	domain := os.Getenv("DOMAIN")
	dbName := os.Getenv("DB")
	mongoURI := os.Getenv("MONGO_URI")
	mode := os.Getenv("MODE")
	_ = mode

	storageInst, err := storage.InitStorage(dbName, mongoURI)
	if err != nil {
		log.Println(err)
		return
	}

	err = storageInst.CreateIndex()
	if err != nil {
		log.Println(err)
	}
	bnsIndexer := bnsindexer.Indexer{
		Storage:      storageInst,
		DefaultRPC:   defaultRPC,
		DefaultRPCWS: defaultRPCWS,
		BNSContract:  constants.BNSContract,
	}
	go bnsIndexer.Start()
	gw := gateway.Gateway{
		DefaultRPC: defaultRPC,
		Port:       port,
		Domain:     domain,
		Storage:    storageInst,
	}
	gw.Start()
}
