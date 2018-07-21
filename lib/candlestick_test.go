package lib

import (
	"fmt"
	"net/http"
	"testing"
	time "time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetCandlestick(t *testing.T) {
	now := time.Now()
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	client.EXPECT().request(&clientOption{
		endpoint: publicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatCandlestick, PairBtcJpy, FiveMinutes, now.Format("20060102")),
	}).Return(
		[]byte(`{"success":1,"data":{"candlestick":[{"type":"5min","ohlcv":[["944912","948000","944155","946404","5.2297",1526860800000],["946001","948200","945336","946845","6.7225",1526861100000]]},{"type":"5min","ohlcv":[["944912","948000","944155","946404","5.2297",1526860800000],["946001","948200","945336","946845","6.7225",1526861100000]]}],"timestamp":1526947199659}}`),
		nil,
	)
	api := NewApi(nil)
	api.client = client
	candlestick, err := api.GetCandlestick(PairBtcJpy, FiveMinutes, now)

	assert.Nil(t, err)
	assert.NotNil(t, candlestick)
	assert.Equal(t, 944912.0, candlestick[0].Ohlcv[0].Open)
	assert.Equal(t, 948000.0, candlestick[0].Ohlcv[0].High)
	assert.Equal(t, 944155.0, candlestick[0].Ohlcv[0].Low)
	assert.Equal(t, 946404.0, candlestick[0].Ohlcv[0].Close)
	assert.Equal(t, 5.2297, candlestick[0].Ohlcv[0].Volume)
	assert.Equal(t, time.Unix(1526860800000/1000, 1526860800000%1000*1000000), candlestick[0].Ohlcv[0].Date)

	assert.Equal(t, 946001.0, candlestick[1].Ohlcv[1].Open)
	assert.Equal(t, 948200.0, candlestick[1].Ohlcv[1].High)
	assert.Equal(t, 945336.0, candlestick[1].Ohlcv[1].Low)
	assert.Equal(t, 946845.0, candlestick[1].Ohlcv[1].Close)
	assert.Equal(t, 6.7225, candlestick[1].Ohlcv[1].Volume)
	assert.Equal(t, time.Unix(1526861100000/1000, 1526861100000%1000*1000000), candlestick[1].Ohlcv[1].Date)
}
