package coin_price_bet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var ETX uint32 = 0x03
var EOT uint32 = 0x04
var ENQ uint32 = 0x05

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