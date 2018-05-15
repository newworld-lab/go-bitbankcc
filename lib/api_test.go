package lib

import (
	"testing"
	"time"

	"github.com/newworld-lab/go-bitbankcc/constant"
	"github.com/stretchr/testify/assert"
)

func TestGetTicker(t *testing.T) {
	client := &clientImpl{
		endpoint:       "https://public.bitbank.cc",
		defaultTimeout: time.Duration(1000 * time.Millisecond),
	}
	api := &APIImpl{
		client: client,
	}
	ticker, err := api.GetTicker(constant.PairBtcJpy, nil)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
	assert.NotEqual(t, 0, ticker.Buy)
}

func TestGetDepth(t *testing.T) {
	client := &clientImpl{
		endpoint:       "https://public.bitbank.cc",
		defaultTimeout: time.Duration(5000 * time.Millisecond),
	}
	api := &APIImpl{
		client: client,
	}
	depth, err := api.GetDepth(constant.PairBtcJpy, nil)

	assert.Nil(t, err)
	assert.NotNil(t, depth)
	assert.NotEqual(t, 0, depth.Asks)
	assert.NotEqual(t, 0, depth.Bids)
}

func TestGetTransactions(t *testing.T) {
	client := &clientImpl{
		endpoint:       "https://public.bitbank.cc",
		defaultTimeout: time.Duration(50000 * time.Millisecond),
	}
	api := &APIImpl{
		client: client,
	}

	transactionsAll, err := api.GetTransactions(constant.PairBtcJpy, nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, transactionsAll[0])
	assert.NotEqual(t, 0, transactionsAll[0].TransactionId)
	assert.NotEqual(t, 0, transactionsAll[0].Side)
	assert.NotEqual(t, 0, transactionsAll[0].Price)
	assert.NotEqual(t, 0, transactionsAll[0].Amount)
	assert.NotEqual(t, 0, transactionsAll[0].ExecutedAt)

	now := time.Now().Add(-24 * time.Hour)
	transactions, err := api.GetTransactions(constant.PairBtcJpy, &now, nil)
	assert.Nil(t, err)
	assert.NotNil(t, transactions[0])
	assert.NotEqual(t, 0, transactions[0].TransactionId)
	assert.NotEqual(t, 0, transactions[0].Side)
	assert.NotEqual(t, 0, transactions[0].Price)
	assert.NotEqual(t, 0, transactions[0].Amount)
	assert.NotEqual(t, 0, transactions[0].ExecutedAt)
}
