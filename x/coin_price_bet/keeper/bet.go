package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

func (k Keeper) GetDayInfo(ctx sdk.Context, dayId int64) *types.BetDay {
	dayStoreKey := types.DayInfoStoreKey(dayId)
	store := ctx.KVStore(k.storeKey)
	day := &types.BetDay{}
	if b := store.Get(dayStoreKey); b != nil {
		k.cdc.MustUnmarshalBinaryBare(b, &day)
	}
	return day
}

func (k Keeper) SetDayInfo(ctx sdk.Context, dayId int64, day *types.BetDay) {
	dayStoreKey := types.DayInfoStoreKey(dayId)
	store := ctx.KVStore(k.storeKey)
	store.Set(dayStoreKey, k.cdc.MustMarshalBinaryBare(day))
}

func (k Keeper) GetDayCoinInfo(ctx sdk.Context, dayId int64, coinId int64) *types.BetDayCoin {
	dayCoinStoreKey := types.DayCoinInfoStoreKey(dayId, coinId)
	store := ctx.KVStore(k.storeKey)
	dayCoin := &types.BetDayCoin{}
	if b := store.Get(dayCoinStoreKey); b != nil {
		k.cdc.MustUnmarshalBinaryBare(b, &dayCoin)
	}
	return dayCoin
}

func (k Keeper) SetDayCoinInfo(ctx sdk.Context, dayId int64, coinId int64, dayCoin *types.BetDayCoin) {
	dayCoinStoreKey := types.DayCoinInfoStoreKey(dayId, coinId)
	store := ctx.KVStore(k.storeKey)
	store.Set(dayCoinStoreKey, k.cdc.MustMarshalBinaryBare(dayCoin))
}
