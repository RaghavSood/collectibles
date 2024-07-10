package util

import (
	"fmt"
	"html/template"
	"strings"
	"time"

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

func ShortUTCTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05")
}

func LifespanString(start time.Time, end time.Time) string {
	if end.IsZero() {
		end = time.Now()
	}

	start = start.UTC()
	end = end.UTC()

	return PrettyDuration(end.Sub(start), 3)
}

func PrettyDuration(d time.Duration, significantPlaces int) string {
	const (
		secondsInMinute = 60
		secondsInHour   = 60 * secondsInMinute
		secondsInDay    = 24 * secondsInHour
		secondsInMonth  = 30 * secondsInDay
		secondsInYear   = 12 * secondsInMonth
	)

	totalSeconds := int(d.Seconds())
	years := totalSeconds / secondsInYear
	months := (totalSeconds % secondsInYear) / secondsInMonth
	days := (totalSeconds % secondsInMonth) / secondsInDay
	hours := (totalSeconds % secondsInDay) / secondsInHour
	minutes := (totalSeconds % secondsInHour) / secondsInMinute
	seconds := totalSeconds % secondsInMinute

	parts := []struct {
		value int
		unit  string
	}{
		{years, "year"},
		{months, "month"},
		{days, "day"},
		{hours, "hour"},
		{minutes, "minute"},
		{seconds, "second"},
	}

	result := ""
	count := 0
	for _, part := range parts {
		if part.value > 0 {
			if count > 0 {
				result += ", "
			}
			result += fmt.Sprintf("%d %s", part.value, part.unit)
			if part.value > 1 {
				result += "s"
			}
			count++
			if significantPlaces > 0 && count >= significantPlaces {
				break
			}
		}
	}

	if result == "" {
		result = "0 seconds"
	}

	return result
}
