package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

// Get last price of a coin
func (k Keeper) GetLatestCoinPrice(ctx sdk.Context, coinId int64) int64 {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.LastCoinPriceStoreKey(coinId)
	if b := store.Get(storeKey); b != nil {
		return types.BytesToInt64(b)
	}
	return 0
}

// Store last price of a coin
func (k Keeper) SetLatestCoinPrice(ctx sdk.Context, coinId int64, price int64) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.LastCoinPriceStoreKey(coinId)
	store.Set(storeKey, types.Int64ToBytes(price))
}

// Get list of all prices of a coin in a particular day
func (k Keeper) GetDayCoinPrices(ctx sdk.Context, dayId int64, coinId int64) []int64 {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.DayCoinPricesStoreKey(dayId, coinId)
	prices := []int64{}
	if b := store.Get(storeKey); b != nil {
		k.cdc.MustUnmarshalBinaryBare(b, &prices)
	}
	return prices
}

// Set prices of a coin in a particular day
func (k Keeper) SetDayCoinPrices(ctx sdk.Context, dayId int64, coinId int64, prices []int64) {
	store := ctx.KVStore(k.storeKey)
	storeKey := types.DayCoinPricesStoreKey(dayId, coinId)
	store.Set(storeKey, k.cdc.MustMarshalBinaryBare(prices))
}

// Append price to existing list of prices of a coin
func (k Keeper) AppendTodayCoinPrice(ctx sdk.Context, coinId int64, price int64) {
	dayId := types.GetDayId(ctx.BlockTime().Unix())
	dayCoinPrices := k.GetDayCoinPrices(ctx, dayId, coinId)
	dayCoinPrices = append(dayCoinPrices, price)
	k.SetDayCoinPrices(ctx, dayId, coinId, dayCoinPrices)
}

// // SetBlock saves the given block to the store without performing any validation.
// func (k Keeper) SetBlockCoinPrice(ctx sdk.Context, blockId int64, coinId int64, price int64) {
// 	types.Logger.Info(fmt.Sprintf("Got price update for coin(%d) price(%d) block(%d)", coinId, price, blockId))
//
// 	store := ctx.KVStore(k.storeKey)
// 	blockStoreId := types.BlockStoreKey(uint64(blockId))
// 	var block types.Block
//
// 	if !store.Has(blockStoreId) {
// 		types.Logger.Error(fmt.Sprintf("Block(%d) does not exists. Creating...", blockId))
//
// 		blockTime := ctx.BlockTime().Unix()
// 		store.Set(types.BlockTimeStoreKey(uint64(blockTime)), Int64ToBytes(blockTime)) // time: id
// 		types.Logger.Info(fmt.Sprintf("Stored new block time %d", blockTime))
//
//     todayId := types.GetDayId(ctx.BlockTime().Unix())
//     dayCoinId := types.GetDayCoinId(todayId, coinId)
//     dayCoinStoreKey := types.DayCoinBlockTimesStoreKey(uint64(dayCoinId))
//
// 		blockTimes := []int64{}
// 		if blockTimesBytes := store.Get(dayCoinStoreKey); blockTimesBytes != nil {
// 			k.cdc.MustUnmarshalBinaryBare(blockTimesBytes, &blockTimes)
// 		}
// 		blockTimes = append(blockTimes, blockTime)
// 		store.Set(dayCoinStoreKey, k.cdc.MustMarshalBinaryBare(blockTimes)) // [time, ...]
// 		types.Logger.Info(fmt.Sprintf("Appended new block time %d", len(blockTimes)))
//
// 		block = types.Block{}
// 		block.Time = blockTime
// 		block.Prices = make([]int64, len(types.GetCoins()))
// 	} else {
// 		k.cdc.MustUnmarshalBinaryBare(store.Get(blockStoreId), &block)
// 	}
//
// 	block.Prices[coinId] = price
//
// 	store.Set(blockStoreId, k.cdc.MustMarshalBinaryBare(block)) // [{time, prices}, ...]
// 	types.Logger.Info(fmt.Sprintf("Stored new block %+v", block))
// }
