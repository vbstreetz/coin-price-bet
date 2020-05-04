package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
	"net/http"
	"os/exec"
)

const FAUCET_STAKE_AMOUNT = 1000000
const BCCLI_PROGRAM_PATH = "/opt/bccli"

func faucetHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars[restPathVarAddress]

		command := fmt.Sprintf("bccli tx send requester %s %dstake --keyring-backend test -y", address, FAUCET_STAKE_AMOUNT)
		types.Logger.Info(command)

		cmd := exec.Command(BCCLI_PROGRAM_PATH, "tx", "send", "requester", address, fmt.Sprintf("%dstake", FAUCET_STAKE_AMOUNT), "--keyring-backend", "test", "-y", "--home", "/opt/.bccli")
		_, err := cmd.Output()
		if err != nil {
		  types.Logger.Error(fmt.Sprintf("%s", err))
		}

		rest.PostProcessResponse(w, cliCtx, []byte("{}"))
	}
}
