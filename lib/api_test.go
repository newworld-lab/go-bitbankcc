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
