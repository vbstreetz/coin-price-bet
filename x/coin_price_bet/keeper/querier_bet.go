package keeper

import (
	// "fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
	"strconv"
)

type Info struct {
	FirstDay uint64 `json:"firstDay"`
	// BetchainTransferChannel string `json:"betchainTransferChannel"`
	// GaiaTransferChannel     string `json:"gaiaTransferChannel"`
}

type DayInfo struct {
	GrandPrizeAmount uint64 `json:"grandPrizeAmount"`
	// AtomPriceUSD     uint64   `json:"atomPriceUSD"`
	CoinsPerf   []int64  `json:"coinsPerf"`
	CoinsVolume []uint64 `json:"coinsVolume"`
	State       uint8    `json:"state"`
}

type MyInfo struct {
	TotalBetsAmount uint64 `json:"totalBetsAmount"`
	TotalWinsAmount uint64 `json:"totalWinsAmount"`
}

type MyDayInfo struct {
	CoinBetTotalAmount []uint64 `json:"coinBetTotalAmount"`
	TotalBetAmount     uint64   `json:"totalBetAmount"`
	TotalWinAmount     uint64   `json:"totalWinAmount"`
	Paid               bool     `json:"paid"`
}

// queryInfo is a query function to get general blockchain info
func queryInfo(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery,
) ([]byte, error) {
	// 	betchainTransferChannel, err := keeper.GetChannel(ctx, types.GAIA_CHAIN_ID, types.TRANSFER_PORT)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	gaiaTransferChannel, err := keeper.GetChannel(ctx, types.VB_CHAIN_ID, types.TRANSFER_PORT)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	ret := &Info{
		FirstDay: uint64(types.GetGenesisBlockTime()),
		// 		BetchainTransferChannel: betchainTransferChannel,
		// 		GaiaTransferChannel:     gaiaTransferChannel,
	}
	return keeper.cdc.MustMarshalJSON(ret), nil
}

// queryMyInfo is a query function to get general info of an address
func queryMyInfo(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, bettor string,
) ([]byte, error) {
	ret := &MyInfo{}

	ret.TotalBetsAmount = uint64(keeper.GetTotalBetsAmount(ctx))
	ret.TotalWinsAmount = uint64(keeper.GetTotalWinsAmount(ctx))

	return keeper.cdc.MustMarshalJSON(ret), nil
}

// queryDayInfo is a query function to get general info of a particular day
func queryDayInfo(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, dayIdString string,
) ([]byte, error) {
	ret := &DayInfo{}

	betDayId, err := strconv.ParseInt(dayIdString, 0, 64)
	if err != nil {
		return nil, err
	}

	ret.State = keeper.GetDayState(ctx, betDayId)

	betDay := keeper.GetDayInfo(ctx, betDayId)

	if betDay == nil {
		ret.CoinsPerf = []int64{}
		ret.CoinsVolume = []uint64{}
	} else {
		ret.GrandPrizeAmount = betDay.GrandPrize

		coins := types.GetCoins()

		for coinId := range coins {
			var perf int64
			var volume uint64

			if betDayCoin := keeper.GetDayCoinInfo(ctx, betDayId, int64(coinId)); betDayCoin != nil {
				volume = uint64(betDayCoin.TotalAmount)
			}
			if prices := keeper.GetDayCoinPrices(ctx, betDayId, int64(coinId)); len(prices) > 1 {
				perf = int64(float64(types.MULTIPLIER) * float64(prices[len(prices)-1]-prices[0]) / float64(prices[0]))
			}

			ret.CoinsPerf = append(ret.CoinsPerf, perf)
			ret.CoinsVolume = append(ret.CoinsVolume, volume)
		}
	}

	// 	coins := types.GetCoins()
	//
	// 	found := false
	// 	var atomCoinId int64
	// 	for i, c := range coins {
	// 		if c == "ATOM" {
	// 			atomCoinId = int64(i)
	// 			found = true
	// 			break
	// 		}
	// 	}
	// 	if !found {
	// 		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, fmt.Sprintf("could not compute ATOM coin id"))
	// 	}

	// ret.AtomPriceUSD = uint64(keeper.GetLatestCoinPrice(ctx, atomCoinId))

	return keeper.cdc.MustMarshalJSON(ret), nil
}

// queryDayMyInfo is a query function to get general info of an address on a particular day
func queryMyDayInfo(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, dayIdString string, bettor string,
) ([]byte, error) {
	ret := &MyDayInfo{
		CoinBetTotalAmount: []uint64{},
	}

	betDayId, err := strconv.ParseInt(dayIdString, 0, 64)
	if err != nil {
		return nil, err
	}

	coins := types.GetCoins()
	winningCoinId := int64(keeper.GetWinningDayCoinId(ctx, betDayId))

	for coinId := range coins {
		amount := uint64(keeper.GetDayCoinBettorAmount(ctx, betDayId, int64(coinId), bettor))
		ret.CoinBetTotalAmount = append(ret.CoinBetTotalAmount, amount)
		ret.TotalBetAmount += amount
	}

	ret.TotalWinAmount = uint64(keeper.GetDayWinAmount(ctx, betDayId, winningCoinId, bettor))
	ret.Paid = keeper.GetDayCoinBettorPaid(ctx, betDayId, winningCoinId, bettor)

	return keeper.cdc.MustMarshalJSON(ret), nil
}
