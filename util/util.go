package util

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/RaghavSood/collectibles/prices"
	"github.com/RaghavSood/collectibles/types"
)

func NoEscapeHTML(str string) template.HTML {
	return template.HTML(str)
}

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

func MultiParam(els ...any) []any {
	return els
}

func ItemPercentString(count, total int) string {
	if total == 0 {
		return "0.00%"
	}

	percentage := float64(count) / float64(total) * 100
	return fmt.Sprintf("%.2f%%", percentage)
}
