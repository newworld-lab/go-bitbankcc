package lib

import (
	"net/http"
	"net/url"
	"testing"

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

	client.EXPECT().request(&clientOption{
		endpoint: privateApiEndpoint,
		method:   http.MethodGet,
		path:     "v1/user/spot/trade_history" + query.Encode(),
		header: http.Header{
			"ACCESS-KEY":       {"key"},
			"ACCESS-NONCE":     {"1532146131702"},
			"ACCESS-SIGNATURE": {"7a80043409e092bc32b16f1fac5b6c5f111ab6b2854c7a5f66f2de4252030a91"},
		},
	}).Return(
		[]byte(`{"success": 0,"data":{"trades":[{"trade_id": 0,"pair": "○○○","order_id": 0,"side": "○○○","type": "○○○","amount": "○○○","price": "○○○","maker_taker": "○○○","fee_amount_base": "○○○","fee_amount_quote": "○○○","executed_at": 0}]}}`),
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
	trades, err := api.GetTrades(entity.PairBtcJpy, 10, 0, 20180701, 20180720, entity.Asc)
	assert.Nil(t, err)
	assert.NotNil(t, trades)
	assert.Equal(t, trades[1].TradeId, 0)
	assert.Equal(t, trades[1].Pair, "○○○")
	assert.Equal(t, trades[1].Side, "○○○")
	assert.Equal(t, trades[1].Type, "○○○")
	assert.Equal(t, trades[1].Amount, "○○○")
	assert.Equal(t, trades[1].Price, "○○○")
	assert.Equal(t, trades[1].MakerTaker, "○○○")
	assert.Equal(t, trades[1].FeeAmountBase, "○○○")
	assert.Equal(t, trades[1].FeeAmountQuote, "○○○")
	assert.Equal(t, trades[1].ExecuteAt, 0)
}
