package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
	"strconv"
	"time"
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
	ret := &types.Info{
		FirstDay: uint64(time.Now().Add(time.Duration(-24) * time.Hour).Unix()),
	}
	return keeper.cdc.MustMarshalJSON(ret), nil
}

// queryMyInfo is a query function to get general info of an address
func queryMyInfo(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, address string,
) ([]byte, error) {
	ret := &types.MyInfo{}
	return keeper.cdc.MustMarshalJSON(ret), nil
}

// queryDayInfo is a query function to get general info of a particular day
func queryDayInfo(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, dayId string,
) ([]byte, error) {
	ret := &types.DayInfo{
		CoinsPerf:   []uint8{},
		CoinsVolume: []uint64{},
		State: 0,
	}
	return keeper.cdc.MustMarshalJSON(ret), nil
}

// queryDayMyInfo is a query function to get general info of an address on a particular day
func queryMyDayInfo(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, dayId string, addressId string,
) ([]byte, error) {
	ret := &types.MyDayInfo{
		CoinBetTotalAmount:     []uint64{},
		CoinPredictedWinAmount: []uint64{},
	}
	return keeper.cdc.MustMarshalJSON(ret), nil
}
