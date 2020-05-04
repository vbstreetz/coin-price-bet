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
		GetCmdReadOrder(storeKey, cdc),
		GetCmdLatestCoinPrices(storeKey, cdc),
	)...)

	return coinPriceBetCmd
}

// GetCmdReadOrder queries order by orderID
func GetCmdReadOrder(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  keeper.QueryOrder,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			orderID := args[0]

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s", queryRoute, keeper.QueryOrder, orderID),
				nil,
			)
			if err != nil {
				fmt.Printf("read request fail - %s \n", orderID)
				return nil
			}

			var order types.Order
			if err := cdc.UnmarshalJSON(res, &order); err != nil {
				return err
			}
			return cliCtx.PrintOutput(order)
		},
	}
}

// GetCmdLatestCoinPrices logs latest prices graph of the given coins
func GetCmdLatestCoinPrices(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  keeper.QueryLatestCoinPrices,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			coinID := args[0]

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s", queryRoute, keeper.QueryLatestCoinPrices, coinID),
				nil,
			)
			if err != nil {
				fmt.Printf("read request fail - %s %s\n", coinID, err)
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
