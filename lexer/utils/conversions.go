package utils

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// StringToNumber converts a string to a numerical value.
func StringToNumber(strNum string) (interface{}, error) {
	if strings.Contains(strNum, ".") {
		return strconv.ParseFloat(strNum, 64)
	}
	return strconv.ParseInt(strNum, 10, 64)
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
		return nil, errors.Wrapf(err, "HexToNumber: failed to parse hex string %s", hexString)
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
