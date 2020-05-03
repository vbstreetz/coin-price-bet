package types

import (
	"encoding/binary"
"time"
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

	// DayCoinBlockTimesStoreKeyPrefix is a prefix for storing block times array [time, ...] of a particular day
	DayCoinBlockTimesStoreKeyPrefix = []byte{0x05}

	// BetDaysStoreKey is a prefix for storing block times array [time, ...]
	BetDaysStoreKeyPrefix = []byte{0x06}

	// DayInfoStoreKey is a prefix for storing day infos dayId => {state, ...}
	DayInfoStoreKeyPrefix = []byte{0x07}

	// DayCoinInfoStoreKey is a prefix for storing day+coin infos: dayId+coinId => {bets, ...}
	DayCoinInfoStoreKeyPrefix = []byte{0x08}
)

func uint64ToBytes(num uint64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, num)
	return result
}

func bytesToUint64(b []byte) int64 {
  return int64(binary.BigEndian.Uint64(b))
}

// These should probably be set in the genesis state

func GetCoins() []string {
	return []string{
		"BTC",
		"ETH",
		"LTC",
		"BAND",
		"ATOM",
		"LINK",
		"XTZ",
	}
}

func GetGenesisBlockTime() int64 {
	return 1585699200 // Wednesday, 1 April 2020 00:00:00 GMT
}

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

// BlockStoreKey is a function to generate key for each block time in store
func DayCoinBlockTimesStoreKey(dayCoinId uint64, ) []byte {
	return append(DayCoinBlockTimesStoreKeyPrefix, uint64ToBytes(dayCoinId)...)
}

//

// Get days since epoch
func GetDayId(blockTime int64) int64 {
	return int64(time.Unix(blockTime, 0).Sub(time.Unix(GetGenesisBlockTime(), 0)).Hours() / 24) // int64 rounds down
}

// GetDayCoinId is a function to generate dayId+coinId
func GetDayCoinId(dayId int64, cointId int64) int64 {
	ret := uint64ToBytes(uint64(dayId))
	return bytesToUint64(append(ret, uint64ToBytes(uint64(cointId))...))
}

//

// GetDayIdStoreKe is a function to generate key for each day info struct
func DayInfoStoreKey(dayId int64) []byte {
	return append(DayInfoStoreKeyPrefix, uint64ToBytes(uint64(dayId))...)
}

// GetDayIdStoreKe is a function to generate key for each block time in store
func DayCoinInfoStoreKey(dayCoinId int64) []byte {
	return append(DayCoinInfoStoreKeyPrefix, uint64ToBytes(uint64(dayCoinId))...)
}
