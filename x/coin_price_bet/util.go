package coin_price_bet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	// "github.com/bandprotocol/bandchain/chain/borsh" // nukes dependency tree: todo - pin version
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
)

var ETX uint32 = 0x03
var EOT uint32 = 0x04
var ENQ uint32 = 0x05

func requestCoinPrice(ctx sdk.Context, keeper Keeper, blockId int64, coinId int64, coin string) error {
	oracleScriptID := oracle.OracleScriptID(2)
	askCount := int64(1)
	minCount := int64(1)

	port := types.ORACLE_DATA_REQUEST_PORT
	channelID, err := keeper.GetChannel(ctx, types.BAND_CHAIN_ID, port)
	if err != nil {
		return sdkerrors.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"not found channel to bandchain",
		)
	}

	sourceChannelEnd, found := keeper.ChannelKeeper.GetChannel(ctx, port, channelID)
	if !found {
		return sdkerrors.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"unknown channel %s port coin_price_bet",
			channelID,
		)
	}
	destinationPort := sourceChannelEnd.Counterparty.PortID
	destinationChannel := sourceChannelEnd.Counterparty.ChannelID

	e := NewEncoder()
	e.EncodeString(coin)
	e.EncodeU64(uint64(types.MULTIPLIER))
	calldata := fmt.Sprintf("%x", e.GetEncodedData())

	// calldata, err := cryptoPrice(coin, types.MULTIPLIER)

	if err != nil {
		return sdkerrors.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"error computing coin price calldata%s",
			err,
		)
	}

	packet := oracle.NewOracleRequestPacketData(
		fmt.Sprintf("%s:%d:%d", types.COMPLETE_COIN_PRICE_UPDATE_ORACLE_PACKET_CLIENT_ID_PREFIX, blockId, coinId),
		oracleScriptID,
		calldata,
		askCount,
		minCount,
	)

	sequence, found := keeper.ChannelKeeper.GetNextSequenceSend(
		ctx, port, channelID,
	)
	if !found {
		return sdkerrors.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"unknown sequence number for channel %s port oracle",
			channelID,
		)
	}

	if err := keeper.ChannelKeeper.SendPacket(
		ctx,
		channel.NewPacket(
			packet.GetBytes(),
			sequence,
			port,
			channelID,
			destinationPort,
			destinationChannel,
			1000000000, // Arbitrarily high timeout for now
		),
	); err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		return err
	} else {
		logger.Info(fmt.Sprintf("Pinged coin(%s)", coin))
	}

	return nil
}

func cryptoPrice(symbol string, multiplier int64) (string, error) {
	// script: 2

	var data = []interface{}{
		uint32(ETX),
		[]byte(symbol),
		uint64(multiplier),
	}

	h, err := buildCalldataHex(data)
	return h, err
}

func openWeatherMap(country string, main_field string, sub_field string, multiplier int64) (string, error) {
	// script: 11

	var data = []interface{}{
		uint32(ENQ),
		[]byte(country),
		uint32(EOT),
		[]byte(main_field),
		uint32(EOT),
		[]byte(sub_field),
		uint64(multiplier),
	}

	h, err := buildCalldataHex(data)
	return h, err
}

func buildCalldataHex(data []interface{}) (string, error) {
	calldata := new(bytes.Buffer)

	for _, v := range data {
		err := binary.Write(calldata, binary.LittleEndian, v)
		if err != nil {
			return "", fmt.Errorf("binary.Write failed:", err)
		}
	}

	return fmt.Sprintf("%x", calldata.Bytes()), nil
}
