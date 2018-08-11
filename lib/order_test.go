package lib

import (
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	entity "github.com/newworld-lab/go-bitbankcc/entity"
)

func TestGetOrder(t *testing.T) {
	key := "key"
	secret := "secret"
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)

	client.EXPECT().request(&clientOption{
		endpoint: privateApiEndpoint,
		method:   http.MethodGet,
		path:     "/v1/user/spot/order?order_id=444134158&pair=btc_jpy",
		header: http.Header{
			"ACCESS-KEY":       {"key"},
			"ACCESS-NONCE":     {"1531539009441"},
			"ACCESS-SIGNATURE": {"57e70c47770a82a6049eba39e73366e28e9d86a1177ccbcb4676b65c3fda3100"},
		},
	}).Return(
		[]byte(`{"success": 1,"data": {"order_id": 444134158,"pair": "btc_jpy","side": "buy","type": "limit","start_amount": "0.00100000","remaining_amount": "0.00100000","executed_amount": "0.00000000","price": "1000.0000","average_price": "0.0000","ordered_at": 1533044000829,"status": "UNFILLED"}}`),
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

	order, err := api.GetOrder(entity.GetOrderParams{
		Pair:    entity.PairBtcJpy,
		OrderID: "444134158",
	})
	assert.Nil(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, 444134158, order.OrderID)
	assert.Equal(t, 0.001, order.StartAmount)
	assert.Equal(t, 1000.0, order.Price)
}

func TestGetOrders(t *testing.T) {
	key := "key"
	secret := "secret"

	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)

	client.EXPECT().request(&clientOption{
		endpoint: privateApiEndpoint,
		method:   http.MethodGet,
		path:     "/v1/user/spot/active_orders?pair=btc_jpy",
		header: http.Header{
			"ACCESS-KEY":       {"key"},
			"ACCESS-NONCE":     {"1531539009441"},
			"ACCESS-SIGNATURE": {"f7ac80894ac0a249fa7845429d864ec151f238e97858339192c3e4cdb75ef2aa"},
		},
	}).Return(
		[]byte(`{"success":1,"data":{"orders":[{"order_id":444134158,"pair":"btc_jpy","side":"buy","type":"limit","start_amount":"0.00100000","remaining_amount":"0.00100000","executed_amount":"0.00000000","price":"1000.0000","average_price":"0.0000","ordered_at":1533044000829,"status":"UNFILLED"}]}}`),
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

	orders, err := api.GetActiveOrders(entity.GetActiveOrdersParams{
		Pair: entity.PairBtcJpy,
	})

	assert.Nil(t, err)
	assert.NotNil(t, orders)
	assert.Equal(t, 444134158, orders[0].OrderID)
	assert.Equal(t, 0.001, orders[0].StartAmount)
	assert.Equal(t, 1000.0, orders[0].Price)
}

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
