package protocol

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func NormalizeAmount(value string, currency string) string {
	coinCurrency := Coin(currency)
	switch coinCurrency {
	case EGLD:
		inputValue, _ := strconv.ParseFloat(value, 64)
		scaleFactor := new(big.Float).SetFloat64(EGLD_DECIMALS)
		amountBigInt := new(big.Float).Mul(new(big.Float).SetFloat64(inputValue), scaleFactor)

		amountBigIntInt, _ := amountBigInt.Int(nil)
		return amountBigIntInt.String()

	case BTC:
		amount, _ := strconv.ParseFloat(value, 64)
		amountBig := big.NewInt(int64(amount * 1e8))
		return amountBig.String()
	default:
		return value
	}
}

func ConvertToFloat(value string, currency string) string {
	coinCurrency := Coin(currency)
	switch coinCurrency {
	case EGLD:
		bigIntVal := new(big.Int)
		bigIntVal.SetString(value, 10)

		floatVal := new(big.Float).SetInt(bigIntVal)
		scaleFactor := new(big.Float).SetFloat64(EGLD_DECIMALS)
		amountVal := new(big.Float).Quo(floatVal, scaleFactor)
		amountStr := fmt.Sprintf("%f", amountVal)
		amountStr = strings.TrimRight(amountStr, "0")
		amountStr = strings.TrimRight(amountStr, ".")

		return amountStr

	case BTC:
		bigIntVal := new(big.Int)
		bigIntVal.SetString(value, 10)

		floatVal := new(big.Float).SetInt(bigIntVal)
		scaleFactor := new(big.Float).SetFloat64(BTC_DECIMALS)
		amountVal := new(big.Float).Quo(floatVal, scaleFactor)
		amountStr := fmt.Sprintf("%f", amountVal)
		amountStr = strings.TrimRight(amountStr, "0")
		amountStr = strings.TrimRight(amountStr, ".")

		return amountStr
	default:
		return value
	}
}