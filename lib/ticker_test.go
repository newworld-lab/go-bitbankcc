package lib

import (
	"fmt"
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetTicker(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	client.EXPECT().request(&clientOption{
		endpoint: publicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatTicker, entity.PairBtcJpy),
	}).Return(
		[]byte(`{"success":1,"data":{"sell":"1020979","buy":"1020712","high":"1023889","low":"963930","last":"1020984","vol":"2075.8257","timestamp":1524573765864}}`),
		nil,
	)
	api := &APIImpl{
		client: client,
	}
	ticker, err := api.GetTicker(entity.PairBtcJpy)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
	assert.Equal(t, 1020712.0, ticker.Buy)
	assert.Equal(t, 1020979.0, ticker.Sell)
	assert.Equal(t, 1023889.0, ticker.High)
	assert.Equal(t, 963930.0, ticker.Low)
	assert.Equal(t, 1020984.0, ticker.Last)
	assert.Equal(t, 2075.8257, ticker.Vol)
}
