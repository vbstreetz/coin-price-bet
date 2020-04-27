package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	BankKeeper    types.BankKeeper
	ChannelKeeper types.ChannelKeeper
}

// NewKeeper creates a new band consumer Keeper instance.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankKeeper types.BankKeeper,
	channelKeeper types.ChannelKeeper,
) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		BankKeeper:    bankKeeper,
		ChannelKeeper: channelKeeper,
	}
}
