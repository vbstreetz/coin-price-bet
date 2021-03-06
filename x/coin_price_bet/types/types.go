package types

import (
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

var Logger log.Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))

var COMPLETE_COIN_PRICE_UPDATE_ORACLE_PACKET_CLIENT_ID_PREFIX string = "COIN_PRICE_UPDATE_REQUEST"

var VB_CHAIN_ID string = "band-consumer"
var GAIA_CHAIN_ID string = "band-cosmoshub"
var BAND_CHAIN_ID string = "ibc-bandchain"

var TRANSFER_PORT string = "transfer"
var ORACLE_PORT string = "coinpricebet"

var MULTIPLIER int64 = 1000000

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
	TotalAmount uint64
	// no amino support for maps
	// 	Bets        map[string]uint64 // address => uint64
	// 	PaidBettors map[string]bool   // address => bool
}

// Structure with all the current bets information in a contest period (e.g. day)
type BetDay struct {
	GrandPrize uint64 // total prize for a day
}
