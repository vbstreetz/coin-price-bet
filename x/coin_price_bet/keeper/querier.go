package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

const (
	QueryOrder           = "order"
	QueryTodayCoinPrices = "today-coin-prices"
	QueryInfo            = "info"
	QueryDayInfo         = "day-info"
)

// NewQuerier is the module level router for state queries.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		types.Logger.Info(fmt.Sprintf("query /%s", path[0]))
		switch path[0] {
		case QueryOrder:
			if len(path) == 1 {
				return queryOrder(ctx, keeper, req, path[1])
			}
		case QueryTodayCoinPrices:
			if len(path[1:]) == 1 {
				return queryTodayCoinPrices(ctx, keeper, req, path[1])
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
