package eth

import (
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Provider struct {
	Client *ethclient.Client
}

func GetProvider(url string) *Provider {
	client, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}
	return &Provider{client}
}

func (p Provider) GetBalance(address string, blockNumber int64) *big.Int {
	addr := common.HexToAddress(address)

	var bigInt *big.Int
	if blockNumber != 0 {
		bigInt = big.NewInt(blockNumber)
	}

	bal, err := p.Client.BalanceAt(
		context.Background(),
		addr, bigInt,
	)
	if err != nil {
		panic(err)
	}

	return bal
}

