package util

import (
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
