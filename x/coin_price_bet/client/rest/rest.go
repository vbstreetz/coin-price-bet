package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

const (
	restPathVarOrderId = "order_id"
	restPathVarCoinId  = "coin_id"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/buy", storeName), buyGoldRequestHandler(cliCtx, storeName)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/order/{%s}", storeName, restPathVarOrderId), readOrderHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/latest-coin-prices/{%s}", storeName, restPathVarCoinId), getLatestCoinPricesHandler(cliCtx, storeName)).Methods("GET")
}
