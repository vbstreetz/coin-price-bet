package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

// SetBlock saves the given block to the store without performing any validation.
func (k Keeper) SetBlockCoinPrice(ctx sdk.Context, blockId int64, coinId int64, price int64) {
	types.Logger.Info(fmt.Sprintf("Got price update for coin(%d) price(%d) block(%d)", coinId, price, blockId))

	store := ctx.KVStore(k.storeKey)
	blockStoreId := types.BlockStoreKey(uint64(blockId))
	var block types.Block

	if !store.Has(blockStoreId) {
		types.Logger.Error(fmt.Sprintf("Block(%d) does not exists. Creating...", blockId))

		blockTime := ctx.BlockTime().Unix()
		b := new(bytes.Buffer)
		binary.Write(b, binary.LittleEndian, int64(blockId))
		store.Set(types.BlockTimeStoreKey(uint64(blockTime)), b.Bytes()) // time: id
		types.Logger.Info(fmt.Sprintf("Stored new block time %d", blockTime))

    todayId := types.GetDayId(ctx.BlockTime().Unix())
    dayCoinId := types.GetDayCoinId(todayId, coinId)
    dayCoinStoreKey := types.DayCoinBlockTimesStoreKey(uint64(dayCoinId))

		blockTimes := []int64{}
		if blockTimesBytes := store.Get(dayCoinStoreKey); blockTimesBytes != nil {
			k.cdc.MustUnmarshalBinaryBare(blockTimesBytes, &blockTimes)
		}
		blockTimes = append(blockTimes, blockTime)
		store.Set(dayCoinStoreKey, k.cdc.MustMarshalBinaryBare(blockTimes)) // [time, ...]
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

// Get prices for the last 3 days
func (k Keeper) GetLatestCoinPriceGraph(ctx sdk.Context, coinId uint64) (*types.PriceGraph, error) {
	store := ctx.KVStore(k.storeKey)

    todayId := types.GetDayId(ctx.BlockTime().Unix())
    dayCoinId := types.GetDayCoinId(todayId, int64(coinId))
    dayCoinStoreKey := types.DayCoinBlockTimesStoreKey(uint64(dayCoinId))

  // todo, compute for last 3 days
  blockTimes := []int64{}
  if blockTimesBytes := store.Get(dayCoinStoreKey); blockTimesBytes != nil {
    k.cdc.MustUnmarshalBinaryBare(blockTimesBytes, &blockTimes)
  }

	graph := &types.PriceGraph{}

	n := 50
	if len(blockTimes) > n {
		blockTimes = blockTimes[len(blockTimes)-n:]
	}

	for _, blockTime := range blockTimes {
		var blockId int64
		if blockIdBytes := store.Get(types.BlockTimeStoreKey(uint64(blockTime))); blockIdBytes == nil {
			types.Logger.Error(fmt.Sprintf("no block id has been recorded for time(%d)", blockTime))
			continue
		} else if err := binary.Read(bytes.NewReader(blockIdBytes), binary.LittleEndian, &blockId); err != nil {
			types.Logger.Error(fmt.Sprintf("%x could not be decoded to int", blockIdBytes))
			continue
		}

		// types.Logger.Info(fmt.Sprintf("Graph for block %d", blockId))

		block := types.Block{}
		if blockBytes := store.Get(types.BlockStoreKey(uint64(blockId))); blockBytes == nil {
			types.Logger.Error(fmt.Sprintf("no block info has been recorded for id(%d)", blockId))
			continue
		} else {
			k.cdc.MustUnmarshalBinaryBare(blockBytes, &block)
		}

		graph.Times = append(graph.Times, block.Time)
		graph.Prices = append(graph.Prices, block.Prices[coinId])
	}

	return graph, nil
}
