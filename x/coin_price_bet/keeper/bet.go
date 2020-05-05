package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

func (k Keeper) GetDayInfo(ctx sdk.Context, dayId int64) *types.BetDay {
	storeKey := types.DayInfoStoreKey(dayId)
	store := ctx.KVStore(k.storeKey)
	day := &types.BetDay{}
	if b := store.Get(storeKey); b != nil {
		k.cdc.MustUnmarshalBinaryBare(b, &day)
	}
	return day
}

func (k Keeper) SetDayInfo(ctx sdk.Context, dayId int64, day *types.BetDay) {
	storeKey := types.DayInfoStoreKey(dayId)
	store := ctx.KVStore(k.storeKey)
	store.Set(storeKey, k.cdc.MustMarshalBinaryBare(*day))
}

func (k Keeper) GetDayCoinInfo(ctx sdk.Context, dayId int64, coinId int64) *types.BetDayCoin {
	storeKey := types.DayCoinInfoStoreKey(dayId, coinId)
	store := ctx.KVStore(k.storeKey)
	dayCoin := &types.BetDayCoin{}
	if b := store.Get(storeKey); b != nil {
		k.cdc.MustUnmarshalBinaryBare(b, &dayCoin)
	}
	return dayCoin
}

func (k Keeper) SetDayCoinInfo(ctx sdk.Context, dayId int64, coinId int64, dayCoin *types.BetDayCoin) {
	storeKey := types.DayCoinInfoStoreKey(dayId, coinId)
	store := ctx.KVStore(k.storeKey)
	types.Logger.Error(fmt.Sprintf("%+v", dayCoin))
	if b, err := k.cdc.MarshalBinaryBare(dayCoin); err != nil {
		types.Logger.Error(fmt.Sprintf("err"))
		types.Logger.Error(fmt.Sprintf("%s", err))
	} else {
		types.Logger.Error(fmt.Sprintf("storing %x", b))
		store.Set(storeKey, b)
	}
}

func (k Keeper) GetDayCoinBettorAmount(ctx sdk.Context, dayId int64, coinId int64, bettor string) int64 {
	storeKey := types.DayCoinBettorAmountStoreKey(dayId, coinId, bettor)
	store := ctx.KVStore(k.storeKey)
	if b := store.Get(storeKey); b != nil {
		return types.BytesToInt64(b)
	} else {
		return 0
	}
}

func (k Keeper) SetDayCoinBettorAmount(ctx sdk.Context, dayId int64, coinId int64, bettor string, amount int64) {
	storeKey := types.DayCoinBettorAmountStoreKey(dayId, coinId, bettor)
	store := ctx.KVStore(k.storeKey)
	store.Set(storeKey, types.Int64ToBytes(amount))
}

func (k Keeper) GetTotalBetsAmount(ctx sdk.Context) int64 {
	storeKey := types.TotalBetsAmountStoreKey
	store := ctx.KVStore(k.storeKey)
	if b := store.Get(storeKey); b != nil {
		return types.BytesToInt64(b)
	} else {
		return 0
	}
}

func (k Keeper) SetTotalBetsAmount(ctx sdk.Context, amount int64) {
	storeKey := types.TotalBetsAmountStoreKey
	store := ctx.KVStore(k.storeKey)
	store.Set(storeKey, types.Int64ToBytes(amount))
}

func (k Keeper) GetTotalWinsAmount(ctx sdk.Context) int64 {
	storeKey := types.TotalWinsAmountStoreKey
	store := ctx.KVStore(k.storeKey)
	if b := store.Get(storeKey); b != nil {
		return types.BytesToInt64(b)
	} else {
		return 0
	}
}

func (k Keeper) SetTotalWinsAmount(ctx sdk.Context, amount int64) {
	storeKey := types.TotalWinsAmountStoreKey
	store := ctx.KVStore(k.storeKey)
	store.Set(storeKey, types.Int64ToBytes(amount))
}
