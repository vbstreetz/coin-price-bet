package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RouterKey is they name of the coin_price_bet module
const RouterKey = ModuleName

// MsgSetSourceChannel is a message for setting source channel to other chain
type MsgSetSourceChannel struct {
	ChainName     string         `json:"chain_name"`
	SourcePort    string         `json:"source_port"`
	SourceChannel string         `json:"source_channel"`
	Signer        sdk.AccAddress `json:"signer"`
}

func NewMsgSetSourceChannel(
	chainName, sourcePort, sourceChannel string,
	signer sdk.AccAddress,
) MsgSetSourceChannel {
	return MsgSetSourceChannel{
		ChainName:     chainName,
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		Signer:        signer,
	}
}

// Route implements the sdk.Msg interface for MsgSetSourceChannel.
func (msg MsgSetSourceChannel) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgSetSourceChannel.
func (msg MsgSetSourceChannel) Type() string { return "set_source_channel" }

// ValidateBasic implements the sdk.Msg interface for MsgSetSourceChannel.
func (msg MsgSetSourceChannel) ValidateBasic() error {
	// TODO: Add validate basic
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgSetSourceChannel.
func (msg MsgSetSourceChannel) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// GetSignBytes implements the sdk.Msg interface for MsgSetSourceChannel.
func (msg MsgSetSourceChannel) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgBuyGold is a message for creating order to buy gold
type MsgBuyGold struct {
	Buyer  sdk.AccAddress `json:"buyer"`
	Amount sdk.Coins      `json:"amount"`
}

// NewMsgBuyGold creates a new MsgBuyGold instance.
func NewMsgBuyGold(
	buyer sdk.AccAddress,
	amount sdk.Coins,
) MsgBuyGold {
	return MsgBuyGold{
		Buyer:  buyer,
		Amount: amount,
	}
}

// Route implements the sdk.Msg interface for MsgBuyGold.
func (msg MsgBuyGold) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgBuyGold.
func (msg MsgBuyGold) Type() string { return "buy_gold" }

// ValidateBasic implements the sdk.Msg interface for MsgBuyGold.
func (msg MsgBuyGold) ValidateBasic() error {
	if msg.Buyer.Empty() {
		return sdkerrors.Wrapf(ErrInvalidBasicMsg, "MsgBuyGold: Sender address must not be empty.")
	}
	if msg.Amount.Empty() {
		return sdkerrors.Wrapf(ErrInvalidBasicMsg, "MsgBuyGold: Amount must not be empty.")
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgBuyGold.
func (msg MsgBuyGold) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

// GetSignBytes implements the sdk.Msg interface for MsgBuyGold.
func (msg MsgBuyGold) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgPlaceBet is a message for placing a bet
type MsgPlaceBet struct {
	Bettor sdk.AccAddress `json:"bettor"`
	Amount sdk.Coins      `json:"amount"`
	CoinId uint8          `json:"coinId"`
}

// NewMsgPlaceBet creates a new MsgPlaceBet instance.
func NewMsgPlaceBet(
	bettor sdk.AccAddress,
	amount sdk.Coins,
	coinId uint8,
) MsgPlaceBet {
	return MsgPlaceBet{
		Bettor: bettor,
		Amount: amount,
		CoinId: coinId,
	}
}

// Route implements the sdk.Msg interface for MsgPlaceBet.
func (msg MsgPlaceBet) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgPlaceBet.
func (msg MsgPlaceBet) Type() string { return "place_bet" }

// ValidateBasic implements the sdk.Msg interface for MsgPlaceBet.
func (msg MsgPlaceBet) ValidateBasic() error {
	if msg.Bettor.Empty() {
		return sdkerrors.Wrapf(ErrInvalidBasicMsg, "MsgPlaceBet: Bettor address must not be empty.")
	}
	if msg.Amount.Empty() {
		return sdkerrors.Wrapf(ErrInvalidBasicMsg, "MsgPlaceBet: Amount must not be empty.")
	}
	if msg.CoinId > uint8(len(GetCoins())) {
		return sdkerrors.Wrapf(ErrInvalidBasicMsg, "MsgPlaceBet: Unknown CoinId.")
	}
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgPlaceBet.
func (msg MsgPlaceBet) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Bettor}
}

// GetSignBytes implements the sdk.Msg interface for MsgPlaceBet.
func (msg MsgPlaceBet) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgPayout is a message for placing a bet
type MsgPayout struct {
	Bettor sdk.AccAddress `json:"bettor"`
	DayId  int64          `json:"dayId"`
}

// NewMsgPayout creates a new MsgPayout instance.
func NewMsgPayout(
	bettor sdk.AccAddress,
	dayId int64,
) MsgPayout {
	return MsgPayout{
		Bettor: bettor,
		DayId:  dayId,
	}
}

// Route implements the sdk.Msg interface for MsgPayout.
func (msg MsgPayout) Route() string { return RouterKey }

// Type implements the sdk.Msg interface for MsgPayout.
func (msg MsgPayout) Type() string { return "place_bet" }

// ValidateBasic implements the sdk.Msg interface for MsgPayout.
func (msg MsgPayout) ValidateBasic() error {
	if msg.Bettor.Empty() {
		return sdkerrors.Wrapf(ErrInvalidBasicMsg, "MsgPayout: Bettor address must not be empty.")
	}
	// if msg.DayId != 0 {
	// 	return sdkerrors.Wrapf(ErrInvalidBasicMsg, "MsgPayout: Unknown DayId.")
	// }
	return nil
}

// GetSigners implements the sdk.Msg interface for MsgPayout.
func (msg MsgPayout) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Bettor}
}

// GetSignBytes implements the sdk.Msg interface for MsgPayout.
func (msg MsgPayout) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
