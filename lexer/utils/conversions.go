package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// StringToNumber converts a string to a numerical value.
func StringToNumber(strNum string) (interface{}, error) {
	if strings.Contains(strNum, ".") {
		return strconv.ParseFloat(strNum, 64)
	}
	return strconv.ParseInt(strNum, 10, 64)
}

func BinaryStringToNumber(binaryNumber string) (interface{}, error) {
	// Convert the binary string to an unsigned integer
	result, err := strconv.ParseUint(binaryNumber, 2, 64)
	if err != nil {
		return nil, fmt.Errorf("BinaryStringToNumber: failed to parse binary string [%w]", err)
	}

	// Check if the result can fit in the smallest uint type
	if result <= math.MaxUint8 {
		return uint8(result), nil
	} else if result <= math.MaxUint16 {
		return uint16(result), nil
	} else if result <= math.MaxUint32 {
		return uint32(result), nil
	} else {
		return uint64(result), nil
	}
}

func HexToNumber(hexString string) (interface{}, error) {
	// Remove the "0x" prefix
	if len(hexString) > 2 && hexString[:2] == "0x" {
		hexString = hexString[2:]
	}

	bitSize := 64

	l := len(hexString)
	if l > 0 && l <= 2 {
		bitSize = 8
	} else if l > 2 && l <= 4 {
		bitSize = 16
	} else if l > 4 && l <= 8 {
		bitSize = 32
	} else if l > 8 {
		bitSize = 64
	}

	// Parse the hexadecimal string to int64
	value, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		return nil, fmt.Errorf("HexToNumber: failed to parse hex string [%w]", err)
	}

	switch bitSize {
	case 8:
		return uint8(value), nil
	case 16:
		return uint16(value), nil
	case 32:
		return uint32(value), nil
	default:
		return uint64(value), nil
	}
}
