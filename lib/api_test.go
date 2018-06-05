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
		endpoint: constant.PublicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatTicker, constant.PairBtcJpy),
	}).Return(
		[]byte(`{"success":1,"data":{"sell":"1020979","buy":"1020712","high":"1023889","low":"963930","last":"1020984","vol":"2075.8257","timestamp":1524573765864}}`),
		nil,
	)
	api := &APIImpl{
		client: client,
	}
	ticker, err := api.GetTicker(constant.PairBtcJpy)
	assert.Nil(t, err)
	assert.NotNil(t, ticker)
	assert.Equal(t, 1020712.0, ticker.Buy)
	assert.Equal(t, 1020979.0, ticker.Sell)
	assert.Equal(t, 1023889.0, ticker.High)
	assert.Equal(t, 963930.0, ticker.Low)
	assert.Equal(t, 1020984.0, ticker.Last)
	assert.Equal(t, 2075.8257, ticker.Vol)
}

func TestGetDepth(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	client.EXPECT().request(&clientOption{
		endpoint: constant.PublicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatDepth, constant.PairBtcJpy),
	}).Return(
		[]byte(`{"success":1,"data":{"asks":[["964745","0.0004"],["964753","0.0010"]],"bids":[["964254","0.0060"],["964249","0.0200"]],"timestamp":1526387708186}}`),
		nil,
	)
	api := &APIImpl{
		client: client,
	}
	depth, err := api.GetDepth(constant.PairBtcJpy)

	assert.Nil(t, err)
	assert.NotNil(t, depth)
	assert.Equal(t, []float64{964745.0, 0.0004}, depth.Asks[0])
	assert.Equal(t, []float64{964254.0, 0.0060}, depth.Bids[0])
}

func TestGetTransactions(t *testing.T) {
	now := time.Now()
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	client.EXPECT().request(&clientOption{
		endpoint: constant.PublicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatTransactionsAll, constant.PairBtcJpy),
	}).Return(
		[]byte(`{"success":1,"data":{"transactions":[{"transaction_id":10654992,"side":"sell","price":"962005","amount":"0.0100","executed_at":1526390858690},{"transaction_id":10654991,"side":"sell","price":"962005","amount":"0.0090","executed_at":1526390856716},{"transaction_id":10654990,"side":"sell","price":"962010","amount":"0.0060","executed_at":1526390856658},{"transaction_id":10654989,"side":"sell","price":"962020","amount":"0.0050","executed_at":1526390856596}]}}`),
		nil,
	)
	client.EXPECT().request(&clientOption{
		endpoint: constant.PublicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatTransactions, constant.PairBtcJpy, now.Format("20060102")),
	}).Return(
		[]byte(`{"success":1,"data":{"transactions":[{"transaction_id":10507587,"side":"sell","price":"956290","amount":"0.0500","executed_at":1526256003017},{"transaction_id":10507588,"side":"sell","price":"956290","amount":"0.0453","executed_at":1526256003458}]}}`),
		nil,
	)

	api := &APIImpl{
		client: client,
	}

	transactionsAll, err := api.GetTransactions(constant.PairBtcJpy, nil)
	assert.Nil(t, err)
	assert.NotNil(t, transactionsAll[0])
	assert.NotEqual(t, 0, transactionsAll[0].TransactionId)
	assert.NotEqual(t, 0, transactionsAll[0].Side)
	assert.NotEqual(t, 0, transactionsAll[0].Price)
	assert.NotEqual(t, 0, transactionsAll[0].Amount)
	assert.NotEqual(t, 0, transactionsAll[0].ExecutedAt)

	transactions, err := api.GetTransactions(constant.PairBtcJpy, &now)
	assert.Nil(t, err)
	assert.NotNil(t, transactions[0])
	assert.NotEqual(t, 0, transactions[0].TransactionId)
	assert.NotEqual(t, 0, transactions[0].Side)
	assert.NotEqual(t, 0, transactions[0].Price)
	assert.NotEqual(t, 0, transactions[0].Amount)
	assert.NotEqual(t, 0, transactions[0].ExecutedAt)
}

func TestGetCandlestick(t *testing.T) {
	now := time.Now()
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	client.EXPECT().request(&clientOption{
		endpoint: constant.PublicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatCandlestick, constant.PairBtcJpy, constant.FiveMinutes, now.Format("20060102")),
	}).Return(
		[]byte(`{"success":1,"data":{"candlestick":[{"type":"5min","ohlcv":[["944912","948000","944155","946404","5.2297",1526860800000],["946001","948200","945336","946845","6.7225",1526861100000]]},{"type":"5min","ohlcv":[["944912","948000","944155","946404","5.2297",1526860800000],["946001","948200","945336","946845","6.7225",1526861100000]]}],"timestamp":1526947199659}}`),
		nil,
	)
	api := NewApi(nil)
	api.client = client
	candlestick, err := api.GetCandlestick(constant.PairBtcJpy, constant.FiveMinutes, now)

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
