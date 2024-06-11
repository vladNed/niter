package mvx

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil/bech32"
)


func ParseDeployResult(result *SmartContractResult) (string, error) {
	resultParts := strings.Split(result.Data, MVX_SEPARATOR)
	if len(resultParts) != 3 {
		return "", errors.New("invalid result data")
	}

	addrBech32, err := Bech32FromHex(resultParts[2])
	if err != nil {
		return addrBech32, nil
	}

	return addrBech32, nil
}

// Converts an mvx from hex to bech32 string format
func Bech32FromHex(addr string) (string, error) {
	if len(addr) != 64 {
		return "", fmt.Errorf("invalid address length %d", len(addr))
	}

	addrDecoded, err := hex.DecodeString(addr)
	if err != nil {
		return "", err
	}

	conv, err := bech32.ConvertBits(addrDecoded, 8, 5, true)
	if err != nil {
		return "", err
	}

	return bech32.Encode(MVX_HRP, conv)
}
