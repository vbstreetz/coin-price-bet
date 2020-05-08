package coin_price_bet

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"strconv"
	"strings"
	// "github.com/bandprotocol/bandchain/chain/borsh" // nukes dependency tree: todo - pin version
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

// NewHandler creates the msg handler of this module, as required by Cosmos-SDK standard.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		// types.Logger.Info(fmt.Sprintf("msg: %+v", msg))
		switch msg := msg.(type) {
		case MsgBuyGold:
			return handleBuyGold(ctx, msg, keeper)
		case MsgSetSourceChannel:
			return handleSetSourceChannel(ctx, msg, keeper)
		case MsgPlaceBet:
			return handlePlaceBet(ctx, msg, keeper)
		case MsgPayout:
			return handlePayout(ctx, msg, keeper)
		// todo: validate msg is from bandchain
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

func handleSetSourceChannel(ctx sdk.Context, msg MsgSetSourceChannel, keeper Keeper) (*sdk.Result, error) {
	keeper.SetChannel(ctx, msg.ChainName, msg.SourcePort, msg.SourceChannel)
	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handleBuyGold(ctx sdk.Context, msg MsgBuyGold, keeper Keeper) (*sdk.Result, error) {
	if err := keeper.EscrowBuyGoldCollateral(ctx, msg.Buyer, msg.Amount); err != nil {
		return nil, err
	}

	orderID := keeper.GetNextOrderCount(ctx)
	keeper.SetOrder(ctx, orderID, types.NewOrder(msg.Buyer, msg.Amount))

	oracleScriptID := oracle.OracleScriptID(3)
	calldata := make([]byte, 8)
	binary.LittleEndian.PutUint64(calldata, uint64(types.MULTIPLIER))
	askCount := int64(1)
	minCount := int64(1)

	port := types.ORACLE_PORT
	channelID, err := keeper.GetChannel(ctx, types.BAND_CHAIN_ID, port)

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
	price := int64(result.Px)

	logger.Info(fmt.Sprintf("Coin(%s) price(%d) time(%d)", coin, price, blockId))

	keeper.SetLatestCoinPrice(ctx, coinId, price)
	keeper.AppendTodayCoinPrice(ctx, coinId, price)

	return nil
}

func handlePlaceBet(ctx sdk.Context, msg MsgPlaceBet, keeper Keeper) (*sdk.Result, error) {
	types.Logger.Info("Placing bet")

	if err := keeper.EscrowBetCollateral(ctx, msg.Bettor, msg.Amount); err != nil {
		return nil, err
	}

	amount := msg.Amount[0].Amount.Uint64()
	bettor := msg.Bettor.String()
	coinId := int64(msg.CoinId)

	types.Logger.Info(fmt.Sprintf("Updating mappings due to bettor(%s) amount(%d) coin(%d)", bettor, amount, coinId))

	betDayId := types.GetDayId(ctx.BlockTime().Unix()) + 1 // for tomorrow

	// Upsert bet day+coin mappings
	betDayCoin := keeper.GetDayCoinInfo(ctx, betDayId, coinId)
	betDayCoin.TotalAmount += amount
	keeper.SetDayCoinInfo(ctx, betDayId, coinId, betDayCoin)

	totalAmount := keeper.GetDayCoinBettorAmount(ctx, betDayId, coinId, bettor)
	keeper.SetDayCoinBettorAmount(ctx, betDayId, coinId, bettor, totalAmount+int64(amount))

	// Upsert bet day mappings
	betDay := keeper.GetDayInfo(ctx, betDayId)
	betDay.GrandPrize += amount
	keeper.SetDayInfo(ctx, betDayId, betDay)

	// Upsert totals
	totalBetsAmount := keeper.GetTotalBetsAmount(ctx)
	keeper.SetTotalBetsAmount(ctx, totalBetsAmount+int64(amount))

	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}

func handlePayout(ctx sdk.Context, msg MsgPayout, keeper Keeper) (*sdk.Result, error) {
	types.Logger.Info("Payout request")

	state := keeper.GetDayState(ctx, msg.DayId)
	bettor := msg.Bettor.String()

	// Ensure day state is PAYOUT
	if state != uint8(types.PAYOUT) {
		return nil, sdkerrors.Wrapf(types.Error, "invalid day(%d) state(%d) coinId(%s)", msg.DayId, state)
	}

	// Get coin ranking
	winningCoinId := int64(keeper.GetWinningDayCoinId(ctx, msg.DayId))

	// Ensure user is a winner
	if keeper.GetDayCoinBettorPaid(ctx, msg.DayId, winningCoinId, bettor) {
		return nil, sdkerrors.Wrapf(types.Error, "bettor(%s) already paid", bettor)
	}
	amount := keeper.GetDayCoinBettorAmount(ctx, msg.DayId, winningCoinId, bettor)
	if amount == 0 {
		return nil, sdkerrors.Wrapf(types.Error, "bettor(%s) has no bets in the winning day(%d), coin(%d)", bettor, msg.DayId, winningCoinId)
	}

	// Transfer prize
	day := keeper.GetDayInfo(ctx, msg.DayId)
	dayCoin := keeper.GetDayCoinInfo(ctx, msg.DayId, winningCoinId)
	prize := amount * int64(day.GrandPrize) / int64(dayCoin.TotalAmount)
	keeper.SetDayCoinBettorPaid(ctx, msg.DayId, winningCoinId, bettor, true)
	keeper.Payout(ctx, msg.Bettor, sdk.NewCoins(sdk.NewInt64Coin("stake", prize)))

	return &sdk.Result{Events: ctx.EventManager().Events().ToABCIEvents()}, nil
}
