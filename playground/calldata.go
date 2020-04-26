package main

import (
	// "bytes"
	// "encoding/binary"
	"fmt"
	// "github.com/bandprotocol/bandchain/chain/borsh"
)

var ETX uint32 = 0x03
var EOT uint32 = 0x04
var ENQ uint32 = 0x05

func main() {
	if h, err := cryptoPrice("BTC", 10); err != nil {
		panic(fmt.Sprintf("%s", err))
	} else {
		fmt.Printf("030000004254430a00000000000000 ==\n%s\n", h)
	}

	if h, err := openWeatherMap("paris", "main", "temp", 1); err != nil {
		panic(fmt.Sprintf("%s", err))
	} else {
		fmt.Printf("050000007061726973040000006d61696e0400000074656d700100000000000000 ==\n%s\n", h)
	}
}

func cryptoPrice(symbol string, multiplier int64) (string, error) {
	// script: 2

	// 	var data = []interface{}{
	// 		uint32(ETX),
	// 		[]byte(symbol),
	// 		uint64(multiplier),
	// 	}
	//
	// 	h, err := buildCalldataHex(data)
	// 	return h, err

	e := borsh.NewEncoder()
	e.EncodeString(symbol)
	e.EncodeSigned64(multiplier)
	return fmt.Sprintf("%x", e.GetEncodedData()), nil
}

func openWeatherMap(country string, main_field string, sub_field string, multiplier int64) (string, error) {
	// script: 11

	// 	var data = []interface{}{
	// 		uint32(ENQ),
	// 		[]byte(country),
	// 		uint32(EOT),
	// 		[]byte(main_field),
	// 		uint32(EOT),
	// 		[]byte(sub_field),
	// 		uint64(multiplier),
	// 	}
	//
	// 	h, err := buildCalldataHex(data)
	// 	return h, err

	e := borsh.NewEncoder()
	e.EncodeString(country)
	e.EncodeString(main_field)
	e.EncodeString(sub_field)
	e.EncodeSigned64(multiplier)
	return fmt.Sprintf("%x", e.GetEncodedData()), nil
}

// func buildCalldataHex(data []interface{}) (string, error) {
// 	calldata := new(bytes.Buffer)
//
// 	for _, v := range data {
// 		err := binary.Write(calldata, binary.LittleEndian, v)
// 		if err != nil {
// 			return "", fmt.Errorf("binary.Write failed:", err)
// 		}
// 	}
//
// 	return fmt.Sprintf("%x", calldata.Bytes()), nil
// }
