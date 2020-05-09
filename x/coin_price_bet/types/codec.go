package types

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/cosmos/cosmos-sdk/codec"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitmenttypes "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/types"
)

// ModuleCdc is the codec for the module.
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	channel.RegisterCodec(ModuleCdc)
	commitmenttypes.RegisterCodec(ModuleCdc)
	oracle.RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetSourceChannel{}, "coin_price_bet/SetSourceChannel", nil)
	cdc.RegisterConcrete(MsgBuyGold{}, "coin_price_bet/BuyGold", nil)
	cdc.RegisterConcrete(MsgPlaceBet{}, "coin_price_bet/PlaceBet", nil)
	cdc.RegisterConcrete(MsgPayout{}, "coin_price_bet/Payout", nil)
}
