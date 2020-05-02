package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func readOrderHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars[restPathVarOrderId]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/order/%s", storeName, orderID),
			nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getLatestCoinPricesHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		coinID := vars[restPathVarCoinId]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/latest-coin-prices/%s", storeName, coinID),
			nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getInfoHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/info", storeName),
			nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getMyInfoHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[restPathVarAddress]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/info/%s", storeName, address),
			nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getDayInfoHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		dayId := vars[restPathVarDayId]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/day-info/%s", storeName, dayId),
			nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getMyDayInfoHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		dayId := vars[restPathVarDayId]
		address := vars[restPathVarAddress]

		res, _, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/day-info/%s/%s", storeName, dayId, address),
			nil,
		)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
