package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
	"time"
)

// SetBlock saves the given block to the store without performing any validation.
func (k Keeper) SetBlockCoinPrice(ctx sdk.Context, blockId int64, coinId int64, price int64) {
	types.Logger.Info(fmt.Sprintf("Got price update for coin(%d) price(%d) block(%d)", coinId, price, blockId))

	store := ctx.KVStore(k.storeKey)
	blockStoreId := types.BlockStoreKey(uint64(blockId))
	var block types.Block

	if !store.Has(blockStoreId) {
		types.Logger.Error(fmt.Sprintf("Block(%d) does not exists. Creating...", blockId))

		blockTime := time.Now().Unix()
		b := new(bytes.Buffer)
		binary.Write(b, binary.LittleEndian, int64(blockId))
		store.Set(types.BlockTimeStoreKey(uint64(blockTime)), b.Bytes()) // time: id
		types.Logger.Info(fmt.Sprintf("Stored new block time %d", blockTime))

		blockTimes := []int64{}
		blockTimesBytes := store.Get(types.BlockTimesStoreKey)
		if blockTimesBytes != nil {
			k.cdc.MustUnmarshalBinaryBare(blockTimesBytes, &blockTimes)
		}
		blockTimes = append(blockTimes, blockTime)
		store.Set(types.BlockTimesStoreKey, k.cdc.MustMarshalBinaryBare(blockTimes)) // [time, ...]
		types.Logger.Info(fmt.Sprintf("Appended new block time %d", len(blockTimes)))

		block = types.Block{}
		block.Time = blockTime
		block.Prices = make([]int64, len(types.GetCoins()))
	} else {
		k.cdc.MustUnmarshalBinaryBare(store.Get(blockStoreId), &block)
	}

	block.Prices[coinId] = price

	store.Set(blockStoreId, k.cdc.MustMarshalBinaryBare(block)) // [{time, prices}, ...]
	types.Logger.Info(fmt.Sprintf("Stored new block %+v", block))
}

func (k Keeper) GetLatestCoinPriceGraph(ctx sdk.Context, coinId uint64) (*types.PriceGraph, error) {
	store := ctx.KVStore(k.storeKey)

	blockTimes := []uint64{}
	blockTimesBytes := store.Get(types.BlockTimesStoreKey)
	if blockTimesBytes == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "no times have been recorded")
	}
	k.cdc.MustUnmarshalBinaryBare(blockTimesBytes, &blockTimes)

	graph := &types.PriceGraph{}

	if len(blockTimes) > 10 {
		blockTimes = blockTimes[len(blockTimes)-10:]
	}

	for _, blockTime := range blockTimes {
		var blockId int64
		blockIdBytes := store.Get(types.BlockTimeStoreKey(blockTime))
		if blockIdBytes == nil {
			types.Logger.Error(fmt.Sprintf("no block id has been recorded for time(%d)", blockTime))
			continue
		}
		if err := binary.Read(bytes.NewReader(blockIdBytes), binary.LittleEndian, &blockId); err != nil {
			types.Logger.Error(fmt.Sprintf("%x could not be decoded to int", blockIdBytes))
			continue
		}

		// types.Logger.Info(fmt.Sprintf("Graph for block %d", blockId))

		block := types.Block{}
		blockBytes := store.Get(types.BlockStoreKey(uint64(blockId)))
		if blockBytes == nil {
			types.Logger.Error(fmt.Sprintf("no block info has been recorded for id(%d)", blockId))
			continue
		}
		k.cdc.MustUnmarshalBinaryBare(blockBytes, &block)

		graph.Times = append(graph.Times, block.Time)
		graph.Prices = append(graph.Prices, block.Prices[coinId])
	}

	return graph, nil
}
