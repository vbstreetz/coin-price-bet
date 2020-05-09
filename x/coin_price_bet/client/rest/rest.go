package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

const (
	restPathVarOrderId = "order_id"
	restPathVarCoinId  = "coin_id"
	restPathVarDayId   = "day_id"
	restPathVarAddress = "address"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/buy", storeName), buyGoldRequestHandler(cliCtx, storeName)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/order/{%s}", storeName, restPathVarOrderId), readOrderHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/today-coin-prices/{%s}", storeName, restPathVarCoinId), getTodayCoinPricesHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/place-bet", storeName), placeBetHandler(cliCtx, storeName)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/info", storeName), getInfoHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/info/{%s}", storeName, restPathVarAddress), getMyInfoHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/day-info/{%s}", storeName, restPathVarDayId), getDayInfoHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/day-info/{%s}/{%s}", storeName, restPathVarDayId, restPathVarAddress), getMyDayInfoHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/faucet/{%s}", storeName, restPathVarAddress), faucetHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/payout", storeName), payoutHandler(cliCtx, storeName)).Methods("POST")
}
