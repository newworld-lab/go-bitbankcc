package lib

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/newworld-lab/go-bitbankcc/constant"
	"github.com/stretchr/testify/assert"
)

func TestGetTicker(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	client.EXPECT().request(&clientOption{
		method:  http.MethodGet,
		path:    fmt.Sprintf(formatTicker, constant.PairBtcJpy),
		timeout: time.Duration(0),
	}).Return(
		[]byte(`{"success":1,"data":{"sell":"1020979","buy":"1020712","high":"1023889","low":"963930","last":"1020984","vol":"2075.8257","timestamp":1524573765864}}`),
		nil,
	)
	api := &APIImpl{
		client: client,
	}
	ticker, err := api.GetTicker(constant.PairBtcJpy, nil)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
	assert.Equal(t, float64(1020712), ticker.Buy)
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
