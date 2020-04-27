package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

var Logger log.Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))

var BUY_GOLD_PACKET_CLIENT_ID_PREFIX string = "ORDER:"
var COMPLETE_COIN_PRICE_UPDATE_ORACLE_PACKET_CLIENT_ID_PREFIX string = "COIN_PRICE_UPDATE_REQUEST"

var BANDCHAIN_ID string = "ibc-bandchain"
var BANDCHAIN_PORT string = "coinpricebet"

var MULTIPLIER int64 = 1000000

type OrderStatus uint8

const (
	Pending OrderStatus = iota
	Active
	Completed
)

type Order struct {
	Owner  sdk.AccAddress `json:"owner"`
	Amount sdk.Coins      `json:"amount"`
	Gold   sdk.Coin       `json:"gold"`
	Status OrderStatus    `json:"status"`
}

func NewOrder(owner sdk.AccAddress, amount sdk.Coins) Order {
	return Order{
		Owner:  owner,
		Amount: amount,
		Status: Pending,
	}
}

func GetCoins() []string {
	return []string{
		"BTC",
		"ETH",
		"LTC",
		"BAND",
		"ATOM",
		"LINK",
		"XTZ",
	}
}

type PriceGraph struct {
	Times  []int64 `json:"times"`
	Prices []int64 `json:"prices"`
}

type Block struct {
	Time   int64   `json:"times"`
	Prices []int64 `json:"prices"`
}
