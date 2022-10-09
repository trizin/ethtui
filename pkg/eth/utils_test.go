package eth

import (
	"math/big"
	"testing"
)

func TestGetEthValue(t *testing.T) {
	wei := big.NewInt(1100000000000000000)
	expected := 1.1
	got := GetEthValue(wei)
	if got != expected {
		t.Errorf("GetEthValue() = %v, want %v", got, expected)
	}
}
