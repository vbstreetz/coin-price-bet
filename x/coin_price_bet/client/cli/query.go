package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/keeper"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

// GetQueryCmd returns
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	coinPriceBetCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the coin_price_bet module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	coinPriceBetCmd.AddCommand(flags.GetCommands(
		GetCmdTodayCoinPrices(storeKey, cdc),
	)...)

	return coinPriceBetCmd
}

// Logs today's prices graph of a given coin
func GetCmdTodayCoinPrices(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  keeper.QueryTodayCoinPrices,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			coinId := args[0]

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s", queryRoute, keeper.QueryTodayCoinPrices, coinId),
				nil,
			)
			if err != nil {
				fmt.Printf("read request fail - %s %s\n", coinId, err)
				return nil
			}

			var prices []int64
			if err := cdc.UnmarshalJSON(res, &prices); err != nil {
				return err
			}
			return cliCtx.PrintOutput(prices)
		},
	}
}
