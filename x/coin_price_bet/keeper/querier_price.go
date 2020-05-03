package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

// queryLatestCoinPrices is a query function to get latest `coin` prices
func queryLatestCoinPrices(
	ctx sdk.Context, keeper Keeper, req abci.RequestQuery, coinIdString string,
) ([]byte, error) {
	coinId, err := strconv.ParseInt(coinIdString, 10, 64)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, fmt.Sprintf("wrong format for coinId %s", err.Error()))
	}
	prices := keeper.GetLatestCoinPrices(ctx, coinId)
	return keeper.cdc.MustMarshalJSON(prices), nil
}
