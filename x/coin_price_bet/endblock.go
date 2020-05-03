package coin_price_bet

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

const PRICE_PING_BLOCK_MOD int64 = 300

func handleEndBlock(ctx sdk.Context, k Keeper, block abci.RequestEndBlock) {
	coins := GetCoins()
	coinId := block.Height % PRICE_PING_BLOCK_MOD // Start ping for prices from BTC(0) to last coin
	if int(coinId) < len(coins) {
		blockId := block.Height - coinId
		if err := requestCoinPrice(ctx, k, blockId, coinId, coins[coinId]); err != nil {
			types.Logger.Error(fmt.Sprintf("%s", err))
		}
	}
}

// func handleEndBlock(ctx sdk.Context, k Keeper, block abci.RequestEndBlock) {
// 	if block.Height%10 == 0 {
// 		if err := requestCoinPrice(ctx, k, 0, 0, "BTC"); err != nil {
// 			types.Logger.Error(fmt.Sprintf("%s", err))
// 		}
// 	}
// }
