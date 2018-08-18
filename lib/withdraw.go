package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/pkg/errors"
)

const (
	formatWithdraw        = "/v1/user/withdrawal_account?asset=%s"
	formatRequestWithdraw = "/v1/user/request_withdrawal"
)

type withdrawAccountResponse struct {
	baseResponse
	Data struct {
		baseData
		Accounts accounts `json:"accounts"`
	} `json:"data"`
}

type requestWithdrawResponse struct {
	baseResponse
	Data withdraw `json:"data"`
}

type accounts []account

type account struct {
	entity.Account
}

type withdraw struct {
	baseData
	entity.Withdraw
	RequestedAt int64 `json:"requested_at"`
}

func (as accounts) convert() entity.Accounts {
	accounts := make(entity.Accounts, 0)
	for _, a := range as {
		accounts = append(accounts, entity.Account{
			UUID:    a.Account.UUID,
			Label:   a.Account.Label,
			Address: a.Account.Address,
		})
	}
	return accounts
}

func (w *withdraw) requestConvert() *entity.Withdraw {
	var requestedAt time.Time

	requestedAt = time.Unix(w.RequestedAt/1000, w.RequestedAt%1000*1000000)

	return &entity.Withdraw{
		UUID:        w.UUID,
		Asset:       w.Asset,
		AccountUUID: w.AccountUUID,
		Amount:      w.Amount,
		Fee:         w.Fee,
		Label:       w.Label,
		Address:     w.Address,
		Txid:        w.Txid,
		Status:      w.Status,
		RequestedAt: requestedAt,
	}
}

func (api *APIImpl) GetWithdraw(asset entity.TypeAsset) (entity.Accounts, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	if api.option == nil || api.option.ApiKey == nil || api.option.ApiSecret == nil {
		return nil, errors.New("ApiKey or ApiSecret is nil")
	}

	var path string
	path = fmt.Sprintf(formatWithdraw, asset)

	header, err := api.createCertificationHeader(path)
	if err != nil {
		return nil, err
	}

	bytes, err := api.client.request(&clientOption{
		endpoint: privateApiEndpoint,
		method:   http.MethodGet,
		path:     path,
		header:   header,
	})

	if err != nil {
		return nil, err
	}

	res := new(withdrawAccountResponse)
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Success != 1 {
		return nil, errors.Errorf("api error code=%d", res.Data.Code)
	}

	return res.Data.Accounts.convert(), nil
}

func (api *APIImpl) PostRequestWithdraw(params entity.PostWithdrawParams) (*entity.Withdraw, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	if api.option == nil || api.option.ApiKey == nil || api.option.ApiSecret == nil {
		return nil, errors.New("ApiKey or ApiSecret is nil")
	}
	path := formatRequestWithdraw

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	header, err := api.createCertificationHeader(string(body))
	if err != nil {
		return nil, err
	}

	bytes, err := api.client.request(&clientOption{
		endpoint: privateApiEndpoint,
		method:   http.MethodPost,
		path:     path,
		header:   header,
		body:     body,
	})
	if err != nil {
		return nil, err
	}

	res := new(requestWithdrawResponse)
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}

	if res.Success != 1 {
		return nil, errors.Errorf("api error code=%d", res.Data.Code)
	}

	return res.Data.requestConvert(), nil
}
