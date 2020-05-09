package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
	"strconv"
)

// Query function to get today's `coin` prices
func queryTodayCoinPrices(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, coinIdString string,
) ([]byte, error) {
	coinId, err := strconv.ParseInt(coinIdString, 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, fmt.Sprintf("wrong format for coinId %s", err.Error()))
	}
	todayId := types.GetDayId(ctx.BlockTime().Unix())
	prices := keeper.GetDayCoinPrices(ctx, todayId, coinId)
	return keeper.cdc.MustMarshalJSON(prices), nil
}
