package eth

import (
	"math"
	"math/big"
)

func GetEthValue(wei *big.Int) float64 {
	fbalance := new(big.Float)
	fbalance.SetString(wei.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	val, _ := ethValue.Float64()
	return val
}
