package bnsindexer

import (
	"bnsportal/abi"
	"bnsportal/constants"
	"bnsportal/models"
	"bnsportal/storage"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type Indexer struct {
	DefaultRPC   string
	DefaultRPCWS string
	BNSContract  string
	Storage      *storage.Storage
}

func (indexer *Indexer) Start() {
	fmt.Println("Starting BFS indexer")
	for {
		err := indexer.IndexBNS()
		if err != nil {
			log.Error().Caller().Err(err).Msg("Error while indexing BFS addresses")
			continue
		}
		time.Sleep(10 * time.Second)
	}
}

func (indexer *Indexer) IndexBNS() error {
	fmt.Println("Indexing BFS addresses")
	evmClient, err := ethclient.Dial(indexer.DefaultRPC)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Error while connecting to node")
		return err
	}
	currentHeight, err := evmClient.BlockNumber(context.Background())
	if err != nil {
		log.Error().Caller().Err(err).Msg("Error while getting current block height")
		return err
	}

	state, err := indexer.Storage.GetIndexerState(context.Background(), "bns_indexer")
	if err != nil {
		if err == mongo.ErrNoDocuments {
			state = &models.IndexerState{
				Indexer:          "bns_indexer",
				LastIndexedBlock: constants.BNSContractDeployedBlock,
			}
			err = indexer.Storage.CreateIndexerState(context.Background(), state)
			if err != nil {
				fmt.Println("Failed to create indexer state: ", err)
				return err
			}
		}
	}

	var startBlock uint64
	var endBlock uint64
	startBlock = state.LastIndexedBlock
	if startBlock == currentHeight {
		return nil
	}
	bnsAddress := common.HexToAddress(indexer.BNSContract)
	bnsContract, err := abi.NewBNS(bnsAddress, evmClient)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Error while loading contract")
		return err
	}

	for {
		if endBlock >= currentHeight {
			break
		}
		if currentHeight-startBlock > 1000 {
			endBlock = startBlock + 1000
		} else {
			endBlock = currentHeight
		}
		opts := &bind.FilterOpts{
			Start:   startBlock,
			End:     &endBlock,
			Context: context.Background(),
		}
		iter, err := bnsContract.BNSFilterer.FilterNameRegistered(opts, nil)
		if err != nil {
			fmt.Println("Failed to get events: ", err)
			continue
		}
		for iter.Next() {
			event := iter.Event

			owner, err := bnsContract.OwnerOf(nil, event.Id)
			if err != nil {
				fmt.Println("Failed to get owner BNS address: ", err)
				continue
			}

			err = indexer.Storage.CreateNameInfo(&models.RegisteredNameInfo{
				Name:              string(event.Name),
				ID:                event.Id.String(),
				RegisteredAtBlock: event.Raw.BlockNumber,
				Owner:             strings.ToLower(owner.String()),
			})
			if err != nil {
				fmt.Println("Failed to create BNS address: ", err)
				continue
			}
			log.Info().Caller().Msgf("Indexed BNS address %s", event.Id.String())
		}

		iter2, err := bnsContract.BNSFilterer.FilterTransfer(opts, nil, nil, nil)
		if err != nil {
			fmt.Println("Failed to get events: ", err)
			continue
		}
		for iter2.Next() {
			event := iter2.Event
			nameID := event.TokenId.String()
			newOwner := strings.ToLower(event.To.String())
			nameInfo, err := indexer.Storage.GetNameInfoByID(nameID)
			if err != nil {
				fmt.Println("Failed to get name info: ", err)
				continue
			}
			if nameInfo.Owner == newOwner {
				continue
			}
			nameInfo.Owner = newOwner
			err = indexer.Storage.UpdateNameInfo(nameInfo)
			if err != nil {
				fmt.Println("Failed to update name info: ", err)
				continue
			}
			log.Info().Caller().Msgf("Updated BNS address %s", nameID)
		}
		state.LastIndexedBlock = endBlock
		err = indexer.Storage.UpdateIndexerState(context.Background(), state)
		if err != nil {
			fmt.Println("Failed to update indexer state: ", err)
			continue
		}
	}
	return nil
}
