package eth

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
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

func (p Provider) SendSignedTransaction(signedhash string) (string, error) {
	rawTx, err := hex.DecodeString(signedhash)
	if err != nil {
		return "", err
	}
	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTx, &tx)
	err = p.Client.SendTransaction(context.Background(), tx)
	if err != nil {
		return "", err
	}
	return tx.Hash().String(), nil
}

func (p Provider) GetBalance(address string, blockNumber uint64) *big.Int {
	addr := common.HexToAddress(address)

	var bigInt *big.Int
	if blockNumber != 0 {
		bigInt = big.NewInt(int64(blockNumber))
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

func (p Provider) GetTransactionReceipt(hash string) (*types.Receipt, error) {
	receipt, err := p.Client.TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (p Provider) GetTransactionInfo(hash string) (*types.Transaction, bool, error) {
	tx, pending, err := p.Client.TransactionByHash(context.Background(), common.HexToHash(hash))
	if err != nil {
		return nil, false, err
	}
	return tx, pending, nil
}