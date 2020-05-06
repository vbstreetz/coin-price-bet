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
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		Short:                      "Coin Price Bet transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	coinPriceBetCmd.AddCommand(flags.PostCommands(
		GetCmdBuyGoldRequest(cdc),
		GetCmdSetChannel(cdc),
		GetCmdPlaceBet(cdc),
	)...)

	return coinPriceBetCmd
}

// GetCmdRequest implements the request command handler.
func GetCmdBuyGoldRequest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy [amount]",
		Short: "Make a new order to buy gold",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Make a new order to buy gold.
Example:
$ %s tx coinpricebet buy 1000000dfsbsdfdf/transfer/uatom
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
$ %s tx coinpricebet set-channel bandchain coin_price_bet dbdfgsdfsd
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

// GetCmdRequest implements the request command handler.
func GetCmdPlaceBet(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "place-bet [coin] [amount]",
		Short: "Place a bet",
		Args:  cobra.ExactArgs(2),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Place a bet on a particular coin.
Example:
$ %s tx coinpricebet place-bet btc 1000000stake
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

			coinIsSupported := false
			var coinId uint8

			coin := strings.ToUpper(args[1])
			for i, c := range types.GetCoins() {
				if c == coin {
					coinId = uint8(i)
					coinIsSupported = true
					break
				}
			}

			if !coinIsSupported {
				return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "coin not supported %s", coin)
			}

			msg := types.NewMsgPlaceBet(
				cliCtx.GetFromAddress(),
				amount,
				coinId,
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
