package util

import (
	"strings"

	"github.com/RaghavSood/collectibles/prices"
	"github.com/RaghavSood/collectibles/types"
)

func BTCValueToUSD(satsValue *types.BigInt) float64 {
	price, err := prices.GetBTCUSDPrice()
	if err != nil {
		price = 0
	}
	satsValueFloat, _ := satsValue.BigFloat().Float64()
	return (satsValueFloat / 100000000) * price
}

func FormatNumber(number string) string {
	// Split the number into integer and decimal parts
	parts := strings.Split(number, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = "." + parts[1]
	}

	// Handle negative numbers
	startOffset := 3
	if integerPart[0] == '-' {
		startOffset++
	}

	// Format the integer part with commas
	for outputIndex := len(integerPart); outputIndex > startOffset; {
		outputIndex -= 3
		integerPart = integerPart[:outputIndex] + "," + integerPart[outputIndex:]
	}

	// Concatenate the formatted integer part with the decimal part
	return integerPart + decimalPart
}
