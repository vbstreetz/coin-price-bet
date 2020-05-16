package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"time"
)

const (
	// ModuleName is the name of the module
	ModuleName = "coinpricebet"
	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	// GlobalStoreKeyPrefix is a prefix for versioning store
	GlobalStoreKeyPrefix = []byte{0x00}
	// ChannelStoreKeyPrefix is a prefix for storing channel
	ChannelStoreKeyPrefix = append(GlobalStoreKeyPrefix, []byte("Channel")...)

	// Prefix to store prices of a coin in a particular day
	DayCoinPricesStoreKeyPrefix = append(GlobalStoreKeyPrefix, []byte("DayCoinPrices")...)

	// LastCoinPriceStoreKeyPrefix is a prefix for storing last price of a coin: coinId => price
	LastCoinPriceStoreKeyPrefix = append(GlobalStoreKeyPrefix, []byte("LastCoinPrice")...)

	// BetDaysStoreKey is a prefix for storing block times array [time, ...]
	BetDaysStoreKeyPrefix = append(GlobalStoreKeyPrefix, []byte("BetDays")...)

	// DayInfoStoreKey is a prefix for storing day infos dayId => {state, ...}
	DayInfoStoreKeyPrefix = append(GlobalStoreKeyPrefix, []byte("DayInfo")...)

	// DayCoinInfoStoreKey is a prefix for storing day+coin infos: dayId+coinId => {bets, ...}
	DayCoinInfoStoreKeyPrefix = append(GlobalStoreKeyPrefix, []byte("DayCoinInfo")...)

	// Prefix for day+coin+bettor => total amount
	DayCoinBettorAmountStoreKeyPrefix = append(GlobalStoreKeyPrefix, []byte("DayCoinBettorAmount")...)

	// Prefix for day+coin+bettor => paid
	DayCoinBettorPaidStoreKeyPrefix = append(GlobalStoreKeyPrefix, []byte("DayCoinBettorPaid")...)

	// All bets amount
	TotalBetsAmountStoreKey = append(GlobalStoreKeyPrefix, []byte("TotalBetsAmount")...)

	// All wins amount
	TotalWinsAmountStoreKey = append(GlobalStoreKeyPrefix, []byte("TotalWinsAmount")...)
)

func UInt64ToBytes(num uint64) []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, num)
	return result
}

// func BytesToUint64(b []byte) uint64 {
// 	return binary.BigEndian.Uint64(b)
// }

func Int64ToBytes(num int64) []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, num)
	return b.Bytes()
}

func BytesToInt64(b []byte) int64 {
	var num int64
	if err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &num); err != nil {
		Logger.Error(fmt.Sprintf("%x could not be decoded to int64", b))
	}
	return num
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

func GetFirstDayId() int64 {
	return GetDayId(GetGenesisBlockTime())
}

func GetEscrowAddress() sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte("COLLATERAL")))
}

// Get days since epoch
func GetDayId(blockTime int64) int64 {
	return int64(time.Unix(blockTime, 0).Sub(time.Unix(GetGenesisBlockTime(), 0)).Hours() / 24) // int64 rounds down
}

// Generate dayId+coinId
func GetDayCoinId(dayId int64, coinId int64) int64 {
	return BytesToInt64(append(Int64ToBytes(dayId), Int64ToBytes(coinId)...))
}

//

// Generate key for each verified channel in store
func ChannelStoreKey(chainName, channelPort string) []byte {
	buf := append(ChannelStoreKeyPrefix, []byte(chainName)...)
	buf = append(buf, []byte(channelPort)...)
	return buf
}

//

// Generate key for each day+coin prices in store
func DayCoinPricesStoreKey(dayId int64, coinId int64) []byte {
	ret := DayCoinPricesStoreKeyPrefix
	ret = append(ret, Int64ToBytes(dayId)...)
	ret = append(ret, Int64ToBytes(coinId)...)
	return ret
}

// Generate key for each coin last price store key
func LastCoinPriceStoreKey(coinId int64) []byte {
	return append(LastCoinPriceStoreKeyPrefix, Int64ToBytes(coinId)...)
}

//

// Generate key for each day info struct
func DayInfoStoreKey(dayId int64) []byte {
	return append(DayInfoStoreKeyPrefix, Int64ToBytes(dayId)...)
}

// Generate key for each day+coin info in store
func DayCoinInfoStoreKey(dayId int64, coinId int64) []byte {
	ret := DayCoinInfoStoreKeyPrefix
	ret = append(ret, Int64ToBytes(dayId)...)
	ret = append(ret, Int64ToBytes(coinId)...)
	return ret
}

// Generate key for each day+coin+bettor => total amount in store
func DayCoinBettorAmountStoreKey(dayId int64, coinId int64, bettor string) []byte {
	ret := DayCoinBettorAmountStoreKeyPrefix
	ret = append(ret, Int64ToBytes(dayId)...)
	ret = append(ret, Int64ToBytes(coinId)...)
	ret = append(ret, []byte(bettor)...)
	return ret
}

// Generate key for each day+coin+bettor => paid in store
func DayCoinBettorPaidStoreKey(dayId int64, coinId int64, bettor string) []byte {
	ret := DayCoinBettorPaidStoreKeyPrefix
	ret = append(ret, Int64ToBytes(dayId)...)
	ret = append(ret, Int64ToBytes(coinId)...)
	ret = append(ret, []byte(bettor)...)
	return ret
}
