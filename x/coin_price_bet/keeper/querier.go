package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
	"strconv"
)

const (
	QueryOrder            = "order"
	QueryLatestCoinPrices = "latest-coin-prices"
	QueryInfo             = "info"
	QueryDayInfo          = "day-info"
)

// NewQuerier is the module level router for state queries.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryOrder:
			if len(path) == 1 {
				return queryOrder(ctx, keeper, req, path[1])
			}
		case QueryLatestCoinPrices:
			if len(path) == 1 {
				return queryLatestCoinPrices(ctx, keeper, req, path[1])
			}
		case QueryInfo:
			switch len(path[1:]) {
			case 0:
				return queryInfo(ctx, keeper, req)
			case 1:
				return queryMyInfo(ctx, keeper, req, path[1])
			}
		case QueryDayInfo:
			switch len(path[1:]) {
			case 1:
				return queryDayInfo(ctx, keeper, req, path[1])
			case 2:
				return queryMyDayInfo(ctx, keeper, req, path[1], path[2])
			}
		}

		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown nameservice query endpoint")
	}
}

// queryOrder is a query function to get order by order ID.
func queryOrder(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, orderId string,
) ([]byte, error) {
	id, err := strconv.ParseInt(orderId, 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, fmt.Sprintf("wrong format for requestid %s", err.Error()))
	}
	order, err := keeper.GetOrder(ctx, uint64(id))
	if err != nil {
		return nil, err
	}
	return keeper.cdc.MustMarshalJSON(order), nil
}

// queryLatestCoinPrices is a query function to get latest `coin` prices
func queryLatestCoinPrices(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, coinIdString string,
) ([]byte, error) {
	coinId, err := strconv.ParseInt(coinIdString, 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, fmt.Sprintf("wrong format for requestid %s", err.Error()))
	}
	graph, err := keeper.GetLatestCoinPriceGraph(ctx, uint64(coinId))
	if err != nil {
		return nil, err
	}
	return keeper.cdc.MustMarshalJSON(graph), nil
}

// queryInfo is a query function to get general blockchain info
func queryInfo(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery,
) ([]byte, error) {
	betchainTransferChannel, err := keeper.GetChannel(ctx, types.GAIA_CHAIN_ID, types.TRANSFER_PORT)
	if err != nil {
		return nil, err
	}

	gaiaTransferChannel, err := keeper.GetChannel(ctx, types.VB_CHAIN_ID, types.TRANSFER_PORT)
	if err != nil {
		return nil, err
	}

	ret := &types.Info{
		FirstDay:                uint64(types.GetGenesisBlockTime()),
		BetchainTransferChannel: betchainTransferChannel,
		GaiaTransferChannel:     gaiaTransferChannel,
	}
	return keeper.cdc.MustMarshalJSON(ret), nil
}

// queryMyInfo is a query function to get general info of an address
func queryMyInfo(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, address string,
) ([]byte, error) {
	ret := &types.MyInfo{}

// 	for _, bet := range allBetsBy {
// 		ret.TotalBetsAmount = 0
// 		ret.TotalWinsAmount = 0
// 	}

	return keeper.cdc.MustMarshalJSON(ret), nil
}

// queryDayInfo is a query function to get general info of a particular day
func queryDayInfo(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, dayIdString string,
) ([]byte, error) {
store := ctx.KVStore(keeper.storeKey)
	ret := &types.DayInfo{}

	dayId, err := strconv.ParseInt(dayIdString, 0, 64)
	if err != nil {
		return nil, err
	}

	todayId := types.GetDayId(ctx.BlockTime().Unix())

  // todo: state
	diff := dayId - todayId
	switch true {
	case diff == 1:
		ret.State = uint8(types.BET)
	case diff == 0:
		ret.State = uint8(types.DRAWING)
	case diff > -1:
		ret.State = uint8(types.PAYOUT)
	default:
		ret.State = uint8(types.INVALID)
	}
	//types.Logger.Info(fmt.Sprintf("%d %d %d %d", dayId, todayId, diff, ret.State))

	  betDayId := types.GetDayId(ctx.BlockTime().Unix())
	  betDayStoreKey := types.DayInfoStoreKey(betDayId)

	betDay := &types.BetDay{}
	if b := store.Get(betDayStoreKey); b == nil {
		ret.CoinsPerf = []uint8{}
		ret.CoinsVolume = []uint64{}
	} else {
		keeper.cdc.MustUnmarshalBinaryBare(b, &betDay)

    ret.GrandPrizeAmount = betDay.GrandPrize

		coins := types.GetCoins()

		for coinId := range coins {
      var perf uint8
      var volume uint64

      betDayCoinId := types.GetDayCoinId(betDayId, int64(coinId))
      betDayCoinStoreKey := types.DayCoinInfoStoreKey(betDayCoinId)
      betDayCoin := &types.BetCoinDay{}
      if b := store.Get(betDayCoinStoreKey); b != nil {
        keeper.cdc.MustUnmarshalBinaryBare(b, &betDayCoin)
        volume = betDayCoin.TotalAmount
      }

			ret.CoinsPerf = append(ret.CoinsPerf, perf) // todo
			ret.CoinsVolume = append(ret.CoinsVolume, volume)
		}
	}

  coins := types.GetCoins()

  found := false
  var atomCoinId uint8
  for i, c := range coins {
    if c == "ATOM" {
      atomCoinId = uint8(i)
      found = true
      break
    }
  }
  if !found {
    return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, fmt.Sprintf("could not compute ATOM coin id"))
  }

  dayCoinId := types.GetDayCoinId(todayId, int64(atomCoinId))
  dayCoinStoreKey := types.DayCoinBlockTimesStoreKey(uint64(dayCoinId))

  blockTimes := []int64{}
  if blockTimesBytes := store.Get(dayCoinStoreKey); blockTimesBytes != nil {
    keeper.cdc.MustUnmarshalBinaryBare(blockTimesBytes, &blockTimes)
  }

  ret.AtomPriceCents =   0 // todo

	return keeper.cdc.MustMarshalJSON(ret), nil
}

// queryDayMyInfo is a query function to get general info of an address on a particular day
func queryMyDayInfo(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, dayId string, address string,
) ([]byte, error) {
  store := ctx.KVStore(keeper.storeKey)

	ret := &types.MyDayInfo{
		CoinBetTotalAmount:     []uint64{},
		CoinPredictedWinAmount: []uint64{},
	}

	betDayId := types.GetDayId(ctx.BlockTime().Unix())
	// betDayStoreKey := types.DayInfoStoreKey(betDayId)
  coins := types.GetCoins()

	for coinId := range coins {
		var amount uint64
		var win uint64

		betDayCoinId := types.GetDayCoinId(betDayId, int64(coinId))
		betDayCoinStoreKey := types.DayCoinInfoStoreKey(betDayCoinId)
		betDayCoin := &types.BetCoinDay{}
		if b := store.Get(betDayCoinStoreKey); b != nil {
			keeper.cdc.MustUnmarshalBinaryBare(b, &betDayCoin)
			amount = betDayCoin.Bets[address]

			//       if !betDayCoin.PaidBettors[address] {
			//         win += betAmount
			//       }
		}
		ret.CoinBetTotalAmount = append(ret.CoinBetTotalAmount, amount)
		ret.CoinPredictedWinAmount = append(ret.CoinPredictedWinAmount, win)
		ret.TotalBetAmount += amount
		ret.TotalWinAmount += win // todo
	}

	return keeper.cdc.MustMarshalJSON(ret), nil
}
