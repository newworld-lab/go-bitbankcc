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
	assert.Equal(t, accounts[0].UUID, "3aba1510-8cab-11e8-xxxx-yyyyyyyyyyyy")
	assert.Equal(t, accounts[0].Label, "bitpay")
	assert.Equal(t, accounts[0].Address, "asdfghjkl1234567890")
}

func TestPostRequestWithdraw(t *testing.T) {
	key := "key"
	secret := "secret"
	ctrl := gomock.NewController(t)
	client := NewMockclient(ctrl)

	client.EXPECT().request(&clientOption{
		endpoint: privateApiEndpoint,
		method:   http.MethodPost,
		path:     "/v1/user/request_withdrawal",
		header: http.Header{
			"ACCESS-KEY":       {"key"},
			"ACCESS-NONCE":     {"1531539009441"},
			"ACCESS-SIGNATURE": {"c46ac3162004a9fe13ac4fcb0771a49e10e4ccd99ce723b0c94ac67197637159"},
		},
		body: []byte(`{"asset":"btc","uuid":"12345","amount":"10","opt_token":"tokenxxx","sms_token":""}`),
	}).Return(
		[]byte(`{"success":1,"data":{"uuid":"1234","asset":"btc","amount":"10","account_uuid":"ssss","fee":"1234","status":"ffff","label":"gggg","txid":"hhhh","address":"jjjj","request_at":1526256003458}}`),
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

	withdraw, err := api.PostRequestWithdraw(entity.PostWithdrawParams{
		Asset:    entity.AssetBtc,
		UUID:     "12345",
		Amount:   10,
		OptToken: "tokenxxx",
		SmsToken: "",
	})

	assert.Nil(t, err)
	assert.NotNil(t, withdraw)
	assert.Equal(t, "1234", withdraw.UUID)
	assert.Equal(t, entity.AssetBtc, withdraw.Asset)
	assert.Equal(t, "ssss", withdraw.AccountUUID)
	assert.Equal(t, 10.0, withdraw.Amount)
	assert.Equal(t, 1234.0, withdraw.Fee)
	assert.Equal(t, "gggg", withdraw.Label)
	assert.Equal(t, "jjjj", withdraw.Address)
	assert.Equal(t, "hhhh", withdraw.Txid)
	assert.Equal(t, "ffff", withdraw.Status)
	assert.NotNil(t, withdraw.RequestedAt)
}
