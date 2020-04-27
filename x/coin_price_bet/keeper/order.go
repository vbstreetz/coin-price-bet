package keeper

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer"
	"github.com/vbstreetz/coin-price-bet/x/coin_price_bet/types"
	"strings"
)

// GetOrderCount returns the current number of all orders ever exist.
func (k Keeper) GetOrderCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.OrdersCountStoreKey)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

// GetNextOrderCount increments and returns the current number of orders.
// If the global order count is not set, it initializes it with value 0.
func (k Keeper) GetNextOrderCount(ctx sdk.Context) uint64 {
	orderCount := k.GetOrderCount(ctx)
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(orderCount + 1)
	store.Set(types.OrdersCountStoreKey, bz)
	return orderCount + 1
}

func (k Keeper) AddOrder(ctx sdk.Context, buyer sdk.AccAddress, amount sdk.Coins) (uint64, error) {
	orderID := k.GetNextOrderCount(ctx)
	// TODO: Config chain name
	collateralChain := "band-cosmoshub"

	// TODO: Support only 1 coin
	if len(amount) != 1 {
		return 0, sdkerrors.Wrapf(types.ErrOnlyOneDenomAllowed, "%d denoms included", len(amount))
	}
	channelID, err := k.GetChannel(ctx, collateralChain, "transfer")
	if err != nil {
		return 0, err
	}
	prefix := transfer.GetDenomPrefix("transfer", channelID)
	if !strings.HasPrefix(amount[0].Denom, prefix) {
		return 0, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom was: %s", amount[0].Denom)
	}

	// escrow source tokens. It fails if balance insufficient.
	escrowAddress := types.GetEscrowAddress()
	err = k.BankKeeper.SendCoins(ctx, buyer, escrowAddress, amount)
	if err != nil {
		return 0, err
	}
	k.SetOrder(ctx, orderID, types.NewOrder(buyer, amount))

	return orderID, nil
}

// SetOrder saves the given order to the store without performing any validation.
func (k Keeper) SetOrder(ctx sdk.Context, id uint64, order types.Order) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OrderStoreKey(id), k.cdc.MustMarshalBinaryBare(order))
}

// GetOrder gets the given order from the store
func (k Keeper) GetOrder(ctx sdk.Context, id uint64) (types.Order, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.OrderStoreKey(id)) {
		return types.Order{}, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "order %d not found", id)
	}
	bz := store.Get(types.OrderStoreKey(id))
	var order types.Order
	k.cdc.MustUnmarshalBinaryBare(bz, &order)
	return order, nil
}
