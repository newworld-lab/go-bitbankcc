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
	assert.Equal(t, 1020712.0, ticker.Buy)
	assert.Equal(t, 1020979.0, ticker.Sell)
	assert.Equal(t, 1023889.0, ticker.High)
	assert.Equal(t, 963930.0, ticker.Low)
	assert.Equal(t, 1020984.0, ticker.Last)
	assert.Equal(t, 2075.8257, ticker.Vol)
	assert.Equal(t, 1524573765864, ticker.Timestamp)
}

func TestGetDepth(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	client.EXPECT().request(&clientOption{
		method:  http.MethodGet,
		path:    fmt.Sprintf(formatDepth, constant.PairBtcJpy),
		timeout: time.Duration(0),
	}).Return(
		[]byte(`{"success":1,"data":{"asks":[["964745","0.0004"],["964753","0.0010"]],"bids":[["964254","0.0060"],["964249","0.0200"]],"timestamp":1526387708186}}`),
		nil,
	)
	api := &APIImpl{
		client: client,
	}
	depth, err := api.GetDepth(constant.PairBtcJpy, nil)

	assert.Nil(t, err)
	assert.NotNil(t, depth)
	assert.Equal(t, []float64{964745.0, 0.0004}, depth.Asks[0])
	assert.Equal(t, []float64{964254.0, 0.0060}, depth.Bids[0])
	assert.Equal(t, 1526387708186, depth.Timestamp)
}
