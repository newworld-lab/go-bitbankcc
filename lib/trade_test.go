package lib

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetTrades(t *testing.T) {
	key := "key"
	secret := "secret"
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	url := &url.URL{}
	query := url.Query()
	query.Add("pair", "btc_jpy")
	query.Add("count", "10")
	query.Add("order_id", "1")
	query.Add("since", "990403199000")
	query.Add("end", "990403199000")
	query.Add("order", "asc")

	client.EXPECT().request(&clientOption{
		endpoint: privateApiEndpoint,
		method:   http.MethodGet,
		path:     "/v1/user/spot/trade_history" + "?" + query.Encode(),
		header: http.Header{
			"ACCESS-KEY":       {"key"},
			"ACCESS-NONCE":     {"1532146131702"},
			"ACCESS-SIGNATURE": {"26319129a8187ef9bfb4fcee6b91c8d08937bff5fcc4cd44db128a6d8757249e"},
		},
	}).Return(
		[]byte(`{"success": 1,"data":{"trades":[{"trade_id": 0,"pair": "btc_jpy","order_id":1,"side": "10","type": "○○○","amount": "10","price": "10","maker_taker": "○○○","fee_amount_base": "10","fee_amount_quote": "10","executed_at": 63662457900}]}}`),
		nil,
	)
	api := &APIImpl{
		client: client,
		option: &APIOption{
			ApiKey:    &key,
			ApiSecret: &secret,
		},
		createAccessNonce: func() int {
			return 1532146131702
		},
	}
	since := time.Date(2001, 5, 20, 23, 59, 59, 0, time.UTC)
	end := time.Date(2002, 5, 20, 23, 59, 59, 0, time.UTC)

	tradeParams := entity.TradeParams{
		Pair:    entity.PairBtcJpy,
		Count:   10,
		OrderId: 1,
		Since:   &since,
		End:     &end,
		Order:   entity.Asc,
	}

	trades, err := api.GetTrades(tradeParams)

	fmt.Println("trades: " + fmt.Sprint(trades))
	fmt.Println("err: " + fmt.Sprint(err))
	assert.Nil(t, err)
	assert.NotNil(t, trades)
	assert.Equal(t, trades[0].TradeId, 0)
	assert.Equal(t, trades[0].Pair, "btc_jpy")
	assert.Equal(t, trades[0].OrderId, 1)
	assert.Equal(t, trades[0].Side, float64(10))
	assert.Equal(t, trades[0].Type, "○○○")
	assert.Equal(t, trades[0].Amount, float64(10))
	assert.Equal(t, trades[0].Price, float64(10))
	assert.Equal(t, trades[0].MakerTaker, "○○○")
	assert.Equal(t, trades[0].FeeAmountBase, float64(10))
	assert.Equal(t, trades[0].FeeAmountQuote, float64(10))
	assert.Equal(t, trades[0].ExecutedAt, time.Unix(63662457900/1000, 63662457900%1000*1000000))
}
