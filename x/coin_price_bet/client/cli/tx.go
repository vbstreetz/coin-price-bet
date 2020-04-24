package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/spf13/cobra"

	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	coinPriceBetCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "coin_price_bet transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	coinPriceBetCmd.AddCommand(flags.PostCommands(
		GetCmdRequest(cdc),
		GetCmdSetChannel(cdc),
		GetCmdGetGoldPrice(cdc),
	)...)

	return coinPriceBetCmd
}

// GetCmdRequest implements the request command handler.
func GetCmdRequest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy [amount]",
		Short: "Make a new order to buy gold",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Make a new order to buy gold.
Example:
$ %s tx coin_price_bet buy 1000000dfsbsdfdf/transfer/uatom
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			amount, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}
			msg := types.NewMsgBuyGold(
				cliCtx.GetFromAddress(),
				amount,
			)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdSetChannel implements the set channel command handler.
func GetCmdSetChannel(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-channel [chain-id] [port] [channel-id]",
		Short: "Register a verified channel",
		Args:  cobra.ExactArgs(3),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Register a verified channel.
Example:
$ %s tx coin_price_bet set-cahnnel bandchain coin_price_bet dbdfgsdfsd
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			msg := types.NewMsgSetSourceChannel(
				args[0],
				args[1],
				args[2],
				cliCtx.GetFromAddress(),
			)

			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdGetGoldPrice implements the get gold price command handler.
func GetCmdGetGoldPrice(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-gold-price",
		Short: "Read current gold price from oracle",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Read current gold price from oracle.
Example:
$ %s tx coin_price_bet get-gold-price
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
// 			inBuf := bufio.NewReader(cmd.InOrStdin())
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
//
// 			msg := types.NewMsgSetSourceChannel(
// 				args[0],
// 				args[1],
// 				args[2],
// 				cliCtx.GetFromAddress(),
// 			)
//
// 			err := msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}
//
// 			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
        return nil
		},
	}

	return cmd
}
