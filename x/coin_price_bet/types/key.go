package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	// ModuleName is the name of the module
	ModuleName = "coinpricebet"
	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	// GlobalStoreKeyPrefix is a prefix for global primitive state variable
	GlobalStoreKeyPrefix = []byte{0x00}

	// OrdersCountStoreKey is a key that help getting to current orders count state variable
	OrdersCountStoreKey = append(GlobalStoreKeyPrefix, []byte("OrdersCount")...)

	// ChannelStoreKeyPrefix is a prefix for storing channel
	ChannelStoreKeyPrefix = []byte{0x01}

	// OrderStoreKeyPrefix is a prefix for storing order
	OrderStoreKeyPrefix = []byte{0x02}

	// BlockStoreKeyPrefix is a prefix for storing block id=>info{time, price}
	BlockStoreKeyPrefix = []byte{0x03}

	// BlockTimeStoreKeyPrefix is a prefix for storing block time=>id
	BlockTimeStoreKeyPrefix = []byte{0x04}

	// BlockTimesStoreKey is a prefix for storing block times array [time, ...]
	BlockTimesStoreKey = []byte{0x05}

	// BetDaysStoreKey is a prefix for storing block times array [time, ...]
	BetDaysStoreKey = []byte{0x06}
)

// ChannelStoreKey is a function to generate key for each verified channel in store
func ChannelStoreKey(chainName, channelPort string) []byte {
	buf := append(ChannelStoreKeyPrefix, []byte(chainName)...)
	buf = append(buf, []byte(channelPort)...)
	return buf
}

// OrderStoreKey is a function to generate key for each order in store
func OrderStoreKey(orderID uint64) []byte {
	return append(OrderStoreKeyPrefix, uint64ToBytes(orderID)...)
}

func uint64ToBytes(num uint64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, num)
	return result
}

func GetEscrowAddress() sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte("COLLATERAL")))
}

// BlockStoreKey is a function to generate key for each block in store
func BlockStoreKey(blockID uint64) []byte {
	return append(BlockStoreKeyPrefix, uint64ToBytes(blockID)...)
}

// BlockStoreKey is a function to generate key for each block time in store
func BlockTimeStoreKey(blockTime uint64) []byte {
	return append(BlockTimeStoreKeyPrefix, uint64ToBytes(blockTime)...)
}
