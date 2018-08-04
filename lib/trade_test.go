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
		path:     "v1/user/spot/trade_history" + "?" + query.Encode(),
		header: http.Header{
			"ACCESS-KEY":       {"key"},
			"ACCESS-NONCE":     {"1532146131702"},
			"ACCESS-SIGNATURE": {"7a80043409e092bc32b16f1fac5b6c5f111ab6b2854c7a5f66f2de4252030a91"},
		},
	}).Return(
		[]byte(`{"success": 1,"data":{"trades":[{"trade_id": 0,"pair": "btc_jpy","order_id":1,"side": "○○○","type": "○○○","amount": "○○○","price": "○○○","maker_taker": "○○○","fee_amount_base": "○○○","fee_amount_quote": "○○○","executed_at": 0}]}}`),
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
	assert.Equal(t, trades[0].Side, "○○○")
	assert.Equal(t, trades[0].Type, "○○○")
	assert.Equal(t, trades[0].Amount, "○○○")
	assert.Equal(t, trades[0].Price, "○○○")
	assert.Equal(t, trades[0].MakerTaker, "○○○")
	assert.Equal(t, trades[0].FeeAmountBase, "○○○")
	assert.Equal(t, trades[0].FeeAmountQuote, "○○○")
	assert.Equal(t, trades[0].ExecuteAt, 0)
}
