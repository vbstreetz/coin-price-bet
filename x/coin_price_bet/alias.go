package coin_price_bet

import (
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/keeper"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
	NewQuerier    = keeper.NewQuerier
	GetCoins      = types.GetCoins
	logger        = types.Logger
)

type (
	Keeper              = keeper.Keeper
	MsgSetSourceChannel = types.MsgSetSourceChannel
	MsgPlaceBet         = types.MsgPlaceBet
	MsgPayout           = types.MsgPayout
)
