package lib

import (
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/newworld-lab/go-bitbankcc/entity"

	"github.com/stretchr/testify/assert"
)

func TestGetWithdraw(t *testing.T) {
	key := "key"
	secret := "secret"
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)
	client.EXPECT().request(&clientOption{
		endpoint: privateApiEndpoint,
		method:   http.MethodGet,
		path:     "/v1/user/withdrawal_account?asset=btc",
		header: http.Header{
			"ACCESS-KEY":       {"key"},
			"ACCESS-NONCE":     {"1531539009441"},
			"ACCESS-SIGNATURE": {"696ad283e185308663ddb4f4e51359ac9aee9ada82d967f68d016f21effb30fb"},
		},
	}).Return(
		[]byte(`{"success": 1,"data": {"accounts": [{"uuid": "3aba1510-8cab-11e8-xxxx-yyyyyyyyyyyy","label": "bitpay","address": "asdfghjkl1234567890"}]}}`),
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

	accounts, err := api.GetWithdraw(entity.AssetBtc)
	assert.Nil(t, err)
	assert.NotNil(t, accounts)
	assert.Equal(t, accounts[0].Uuid, "3aba1510-8cab-11e8-xxxx-yyyyyyyyyyyy")
	assert.Equal(t, accounts[0].Label, "bitpay")
	assert.Equal(t, accounts[0].Address, "asdfghjkl1234567890")
}
