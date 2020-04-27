package coin_price_bet

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	// "github.com/bandprotocol/bandchain/chain/borsh" // nukes dependency tree: todo - pin version
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

// NewHandler creates the msg handler of this module, as required by Cosmos-SDK standard.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgBuyGold:
			return handleBuyGold(ctx, msg, keeper)
		case MsgSetSourceChannel:
			return handleSetSourceChannel(ctx, msg, keeper)
		case channeltypes.MsgPacket:
			var responsePacket oracle.OracleResponsePacketData
			if err := types.ModuleCdc.UnmarshalJSON(msg.GetData(), &responsePacket); err == nil {
				return handleOracleRespondPacketData(ctx, responsePacket, keeper)
			}
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal oracle packet data")
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

//

func handleSetSourceChannel(ctx sdk.Context, msg MsgSetSourceChannel, keeper Keeper) (*sdk.Result, error) {
	keeper.SetChannel(ctx, msg.ChainName, msg.SourcePort, msg.SourceChannel)
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

//

func handleBuyGold(ctx sdk.Context, msg MsgBuyGold, keeper Keeper) (*sdk.Result, error) {
	orderID, err := keeper.AddOrder(ctx, msg.Buyer, msg.Amount)
	if err != nil {
		return nil, err
	}
	oracleScriptID := oracle.OracleScriptID(3)
	calldata := make([]byte, 8)
	binary.LittleEndian.PutUint64(calldata, uint64(types.MULTIPLIER))
	askCount := int64(1)
	minCount := int64(1)

	port := types.BANDCHAIN_PORT
	channelID, err := keeper.GetChannel(ctx, types.BANDCHAIN_ID, port)

	if err != nil {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"not found channel to bandchain",
		)
	}
	sourceChannelEnd, found := keeper.ChannelKeeper.GetChannel(ctx, port, channelID)
	if !found {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"unknown channel %s port coin_price_bet",
			channelID,
		)
	}
	destinationPort := sourceChannelEnd.Counterparty.PortID
	destinationChannel := sourceChannelEnd.Counterparty.ChannelID
	sequence, found := keeper.ChannelKeeper.GetNextSequenceSend(
		ctx, port, channelID,
	)
	if !found {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"unknown sequence number for channel %s port oracle",
			channelID,
		)
	}
	packet := oracle.NewOracleRequestPacketData(
		fmt.Sprintf("%s%d", types.BUY_GOLD_PACKET_CLIENT_ID_PREFIX, orderID),
		oracleScriptID,
		hex.EncodeToString(calldata),
		askCount,
		minCount,
	)
	err = keeper.ChannelKeeper.SendPacket(ctx, channel.NewPacket(packet.GetBytes(),
		sequence, port, channelID, destinationPort, destinationChannel,
		1000000000, // Arbitrarily high timeout for now
	))
	if err != nil {
		return nil, err
	}
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

//

func handleOracleRespondPacketData(ctx sdk.Context, packet oracle.OracleResponsePacketData, keeper Keeper) (*sdk.Result, error) {
	logger.Info(packet.ClientID)

	var err error
	switch true {
	case strings.HasPrefix(packet.ClientID, types.BUY_GOLD_PACKET_CLIENT_ID_PREFIX):
		err = handleOracleRespondFulfillOrder(ctx, packet, keeper)
	case strings.HasPrefix(packet.ClientID, types.COMPLETE_COIN_PRICE_UPDATE_ORACLE_PACKET_CLIENT_ID_PREFIX):
		err = handleOracleCompleteCoinPriceUpdate(ctx, packet, keeper)
	default:
		err = sdkerrors.Wrapf(types.ErrUnknownClientID, "unknown client id %s", packet.ClientID)
	}

	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, err
}

func handleOracleRespondFulfillOrder(ctx sdk.Context, packet oracle.OracleResponsePacketData, keeper Keeper) error {
	logger.Info("Fulfilling buy gold order")
	clientID := strings.Split(packet.ClientID, ":")
	if len(clientID) != 2 {
		return sdkerrors.Wrapf(types.ErrUnknownClientID, "unknown %s client id %s", types.BUY_GOLD_PACKET_CLIENT_ID_PREFIX, packet.ClientID)
	}

	orderID, err := strconv.ParseUint(clientID[1], 10, 64)
	if err != nil {
		return err
	}

	rawResult, err := hex.DecodeString(packet.Result)
	if err != nil {
		return err
	}
	result, err := types.DecodeResult(rawResult)
	if err != nil {
		return err
	}

	// Assume multiplier should be 1000000
	order, err := keeper.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}
	// TODO: Calculate collateral percentage
	goldAmount := order.Amount[0].Amount.Int64() / int64(result.Px)
	if goldAmount == 0 {
		escrowAddress := types.GetEscrowAddress()
		err = keeper.BankKeeper.SendCoins(ctx, escrowAddress, order.Owner, order.Amount)
		if err != nil {
			return err
		}
		order.Status = types.Completed
		keeper.SetOrder(ctx, orderID, order)
	} else {
		goldToken := sdk.NewCoin("gold", sdk.NewInt(goldAmount))
		keeper.BankKeeper.AddCoins(ctx, order.Owner, sdk.NewCoins(goldToken))
		order.Gold = goldToken
		order.Status = types.Active
		keeper.SetOrder(ctx, orderID, order)
	}
	return nil
}

func handleOracleCompleteCoinPriceUpdate(ctx sdk.Context, packet oracle.OracleResponsePacketData, keeper Keeper) error {
	logger.Info("Completing coin price update request")
	clientID := strings.Split(packet.ClientID, ":")
	if len(clientID) != 3 {
		return sdkerrors.Wrapf(types.ErrUnknownClientID, "%s: unknown clientId(%s)", packet.ClientID)
	}

	// Get blockId
	blockId, err := strconv.ParseInt(clientID[1], 10, 64)
	if err != nil {
		return err
	}

	// Get coinId
	coinId, err := strconv.ParseInt(clientID[2], 10, 64)
	if err != nil {
		return err
	}

	coins := GetCoins()

	if int(coinId) >= len(coins) {
		return sdkerrors.Wrapf(types.ErrUnknownClientID, "unknown %s coinId(%s)", types.COMPLETE_COIN_PRICE_UPDATE_ORACLE_PACKET_CLIENT_ID_PREFIX, coinId)
	}
	coin := GetCoins()[coinId]

	// Get price
	rawResult, err := hex.DecodeString(packet.Result)
	if err != nil {
		return err
	}
	result, err := types.DecodeResult(rawResult)
	if err != nil {
		return err
	}
	price := result.Px

	logger.Info(fmt.Sprintf("Coin(%s) price(%d) time(%d)", coin, price, blockId))

	keeper.SetBlockCoinPrice(ctx, int64(blockId), coinId, int64(price))

	return nil
}
