package lib

import (
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAssets(t *testing.T) {
	key := "key"
	secret := "secret"
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	client.EXPECT().request(&clientOption{
		endpoint: privateApiEndpoint,
		method:   http.MethodGet,
		path:     "/v1/user/assets",
		header: http.Header{
			"ACCESS-KEY":       {"key"},
			"ACCESS-NONCE":     {"1531539009441"},
			"ACCESS-SIGNATURE": {"878077fcb096a515390c6ee7d8bf238b8774ed01852feec7bdfacabba607d285"},
		},
	}).Return(
		[]byte(`{"success":1,"data":{"assets":[{"asset":"jpy","amount_precision":4,"onhand_amount":"29123.8952","locked_amount":"0.0000","free_amount":"29123.8952","stop_deposit":false,"stop_withdrawal":false,"withdrawal_fee":{"threshold":"30000.0000","under":"540.0000","over":"756.0000"}},{"asset":"btc","amount_precision":8,"onhand_amount":"0.03440000","locked_amount":"0.00000000","free_amount":"0.03440000","stop_deposit":false,"stop_withdrawal":false,"withdrawal_fee":"0.00100000"}]}}`),
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
	assets, err := api.GetAssets()
	assert.Nil(t, err)
	assert.NotNil(t, assets)
	assert.Equal(t, assets[1].Asset, "btc")
	assert.Equal(t, assets[1].AmountPrecision, 8)
	assert.Equal(t, assets[1].OnhandAmount, 0.0344)
	assert.Equal(t, assets[1].LockedAmount, 0.0)
	assert.Equal(t, assets[1].FreeAmount, 0.0344)
	assert.Equal(t, assets[1].WithDrawalFee, "0.00100000")
}
