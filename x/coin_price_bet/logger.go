package coin_price_bet

import (
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

var logger log.Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
