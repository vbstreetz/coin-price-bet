package keeper

import (
	"fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k Keeper) GetDayState(ctx sdk.Context, dayId int64) uint8 {
	firstDayId := types.GetFirstDayId()
	todayId := types.GetDayId(ctx.BlockTime().Unix())
	tomorrowId := todayId + 1

	types.Logger.Info(fmt.Sprintf("first-day(%d) today(%d) tomorrow(%d) day(%d)", firstDayId, todayId, tomorrowId, dayId))

	// day should be within first day <> tomorrow
	if dayId < firstDayId || dayId > tomorrowId {
		return uint8(types.INVALID)
	}

	// day betting window is open
	if dayId == tomorrowId {
		return uint8(types.BET)
	}

	day := k.GetDayInfo(ctx, dayId)

	// waiting for day to end to compute results
	if dayId == todayId && day.GrandPrize != 0 {
		return uint8(types.DRAWING)
	}

	// winners can withdraw winnings
	return uint8(types.PAYOUT)
}

func (k Keeper) EscrowBetCollateral(ctx sdk.Context, buyer sdk.AccAddress, amount sdk.Coins) error {
	// TODO: Support only 1 coin
	if len(amount) != 1 {
		return sdkerrors.Wrapf(types.ErrOnlyOneDenomAllowed, "%d denoms included", len(amount))
	}
	prefix := "stake"
	if !strings.HasPrefix(amount[0].Denom, prefix) {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom was: %s", amount[0].Denom)
	}
	// Escrow source tokens. It fails if balance insufficient.
	escrowAddress := types.GetEscrowAddress()
	return k.BankKeeper.SendCoins(ctx, buyer, escrowAddress, amount)
}

func (k Keeper) Payout(ctx sdk.Context, buyer sdk.AccAddress, amount sdk.Coins) {
	escrowAddress := types.GetEscrowAddress()
	k.BankKeeper.SendCoins(ctx, escrowAddress, buyer, amount)
}

func (k Keeper) GetDayCoinRanking(ctx sdk.Context, dayId int64) []uint8 {
	coins := types.GetCoins()
	perfCoinMap := map[float64]uint8{}
	var perfs []float64
	for coinId := range coins {
		var perf float64
		if prices := k.GetDayCoinPrices(ctx, dayId, int64(coinId)); len(prices) > 1 {
			perf = float64(types.MULTIPLIER) * float64(prices[len(prices)-1]-prices[0]) / float64(prices[0])
		}
		perfCoinMap[perf] = uint8(coinId)
		perfs = append(perfs, perf)

	}
	sort.Float64s(perfs)
	var ranking []uint8
	for _, perf := range perfs {
		ranking = append(ranking, uint8(perfCoinMap[perf]))
	}
	return ranking
}

func (k Keeper) GetWinningDayCoinId(ctx sdk.Context, dayId int64) uint8 {
	ranking := k.GetDayCoinRanking(ctx, dayId)
	size := len(ranking)
	if size != 0 {
		return ranking[size-1]
	}
	return 255 // invalid
}

func (k Keeper) GetDayCoinBettorPaid(ctx sdk.Context, dayId int64, coinId int64, bettor string) bool {
	storeKey := types.DayCoinBettorPaidStoreKey(dayId, coinId, bettor)
	store := ctx.KVStore(k.storeKey)
	var paid bool
	if b := store.Get(storeKey); b != nil {
		k.cdc.MustUnmarshalBinaryBare(b, &paid)
	}
	return paid
}

func (k Keeper) SetDayCoinBettorPaid(ctx sdk.Context, dayId int64, coinId int64, bettor string, paid bool) {
	storeKey := types.DayCoinBettorPaidStoreKey(dayId, coinId, bettor)
	store := ctx.KVStore(k.storeKey)
	store.Set(storeKey, k.cdc.MustMarshalBinaryBare(paid))
}

// Todo: better algo that also takes into account the other coins
func (k Keeper) GetDayWinAmount(ctx sdk.Context, dayId int64, winningCoinId int64, bettor string) int64 {
	amount := k.GetDayCoinBettorAmount(ctx, dayId, winningCoinId, bettor)
	if amount == 0 {
		return 0
	}

	// Transfer prize
	day := k.GetDayInfo(ctx, dayId)
	dayCoin := k.GetDayCoinInfo(ctx, dayId, winningCoinId)
	return (amount * int64(day.GrandPrize)) / int64(dayCoin.TotalAmount)
}
