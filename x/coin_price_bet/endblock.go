package coin_price_bet

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const PRICE_PING_BLOCK_MOD int64 = 300

func handleEndBlock(ctx sdk.Context, k Keeper, block abci.RequestEndBlock) {
	coins := GetCoins()
	coinId := block.Height % PRICE_PING_BLOCK_MOD // Start ping for prices from BTC(0) to last coin
	if int(coinId) < len(coins) {
		blockId := block.Height - coinId
		requestCoinPrice(ctx, k, blockId, coinId, coins[coinId])
	}
}
