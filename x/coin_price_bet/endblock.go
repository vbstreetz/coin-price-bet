package coin_price_bet

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func handleEndBlock(ctx sdk.Context, k Keeper, block abci.RequestEndBlock) {
	logger.Info(fmt.Sprintf("------> height %d", block.Height))
}
