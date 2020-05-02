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

type Block struct {
	Time   int64   `json:"times"`
	Prices []int64 `json:"prices"`
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

//

type PriceGraph struct {
	Times  []int64 `json:"times"`
	Prices []int64 `json:"prices"`
}

type Info struct {
	FirstDay uint64 `json:"firstDay"`
}

type DayInfo struct {
	GrandPrizeAmount uint64   `json:"grandPrizeAmount"`
	AtomPriceCents   uint8    `json:"atomPriceCents"`
	CoinsPerf        []uint8 `json:"coinsPerf"`
	CoinsVolume      []uint64 `json:"coinsVolume"`
	State uint8            `json:"state"`
}

type MyInfo struct {
	TotalBetsAmount uint64 `json:"totalBetsAmount"`
	TotalWinsAmount uint64 `json:"totalWinsAmount"`
}

type MyDayInfo struct {
	CoinBetTotalAmount     []uint64 `json:"coinBetTotalAmount"`
	CoinPredictedWinAmount []uint64 `json:"coinPredictedWinAmount"`
	TotalBetAmount         uint64   `json:"totalBetAmount"`
	TotalWinAmount         uint64   `json:"totalWinAmount"`
}
