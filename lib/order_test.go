package lib

import (
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	entity "github.com/newworld-lab/go-bitbankcc/entity"
)

func TestPostOrder(t *testing.T) {
	key := "key"
	secret := "secret"
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)

	client.EXPECT().request(&clientOption{
		endpoint: privateApiEndpoint,
		method:   http.MethodPost,
		path:     "/v1/user/spot/order",
		header: http.Header{
			"ACCESS-KEY":       {"key"},
			"ACCESS-NONCE":     {"1531539009441"},
			"ACCESS-SIGNATURE": {"271f954492f8b60d20eca7a4b11e6945e4cd1a0cc2caf0787dbc59749de6d5b5"},
		},
		body: []byte(`{"pair":"btc_jpy","amount":"0.001","price":1000,"side":"buy","type":"limit"}`),
	}).Return(
		[]byte(`{"success":1,"data":{"order_id":444112806,"pair":"btc_jpy","side":"buy","type":"limit","start_amount":"0.00100000","remaining_amount":"0.00100000","executed_amount":"0.00000000","price":"1000.0000","average_price":"0.0000","ordered_at":1533041753546,"status":"UNFILLED"}}`),
		nil,
	)

	api := &APIImpl{
		client: client,
		option: &APIOption{
			ApiKey:    &key,
			ApiSecret: &secret,
		},
		createAccessNonce: func() int {
			return 1531539009441
		},
	}

	order, err := api.PostOrder(entity.PostOrderParams{
		Pair:   entity.PairBtcJpy,
		Amount: 0.001,
		Price:  1000,
		Side:   entity.OrderSideBuy,
		Type:   entity.OrderTypeLimit,
	})

	assert.Nil(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, 444112806, order.OrderID)
	assert.Equal(t, 0.001, order.StartAmount)
	assert.Equal(t, 1000.0, order.Price)
}
