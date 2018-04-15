package lib

import (
	"testing"
	"time"

	"github.com/KoteiIto/go-bitbankcc/constant"
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
