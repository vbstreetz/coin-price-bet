package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

var Logger log.Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))

var BUY_GOLD_PACKET_CLIENT_ID_PREFIX string = "ORDER:"
var COMPLETE_COIN_PRICE_UPDATE_ORACLE_PACKET_CLIENT_ID_PREFIX string = "COIN_PRICE_UPDATE_REQUEST"

var VB_CHAIN_ID string = "band-consumer"
var GAIA_CHAIN_ID string = "band-cosmoshub"
var BAND_CHAIN_ID string = "ibc-bandchain"

var TRANSFER_PORT string = "transfer"
var ORACLE_PORT string = "coinpricebet"

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

type DayState uint8

const (
	BET DayState = iota
	DRAWING
	PAYOUT
	INVALID
)

// Structure with coin bets in a contest period (e.g. day)
type BetDayCoin struct {
	TotalAmount uint64            // ordered ranking after result has been resolved
	Bets        map[string]uint64 // address => uint64
	PaidBettors map[string]bool   // address => bool
}

// Structure with all the current bets information in a contest period (e.g. day)
type BetDay struct {
	GrandPrize      uint64            // total prize for a day
	Ranking         []uint8           // ordered ranking after result has been resolved
	Bets            map[string]uint64 // address => uint64
	ResultCompleted bool
}
