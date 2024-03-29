package eth

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

type Provider struct {
	Client *ethclient.Client
}

func GetProvider(url string) (*Provider, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Provider{client}, nil
}

func (p Provider) SendSignedTransaction(signedhash string) (string, error) {
	// replace 0x with "" if exists
	if signedhash[:2] == "0x" {
		signedhash = signedhash[2:]
	}
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

func (p Provider) GetBalance(address string, blockNumber uint64) (*big.Int, error) {
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
		return nil, err
	}

	return bal, nil
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

func (p Provider) GetNonce(address string) (uint64, error) {
	addr := common.HexToAddress(address)
	nonce, err := p.Client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (p Provider) GetBlockInfo(blockNumber uint64) (*types.Block, error) {
	block, err := p.Client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (p Provider) GetChainId() (*big.Int, error) {
	chainId, err := p.Client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	return chainId, nil
}

func (p Provider) GetGasPrice() (*big.Int, error) {
	gasPrice, err := p.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

func (p Provider) GetGasTipCap() (*big.Int, error) {
	gasLimit, err := p.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, err
	}
	return gasLimit, nil
}

func (p Provider) GetEstimatedGasUsage(
	data []byte,
) (uint64, error) {
	a := common.Address{}
	msg := ethereum.CallMsg{
		From: common.Address{},
		To:   &a,
		Data: data,
	}

	gasLimit, err := p.Client.EstimateGas(context.Background(), msg)

	if err != nil {
		return 0, err
	}

	return gasLimit, nil
}
