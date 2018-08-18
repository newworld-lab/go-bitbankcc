package lib

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/pkg/errors"
)

const (
	formatWithdraw = "/v1/user/withdrawal_account?asset=%s"
)

type withdrawResponse struct {
	baseResponse
	Data struct {
		baseData
		Accounts accounts `json:"accounts"`
	} `json:"data"`
}

type accounts []account

type account struct {
	entity.Account
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

	res := new(withdrawResponse)
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Success != 1 {
		return nil, errors.Errorf("api error code=%d", res.Data.Code)
	}

	return res.Data.Accounts.convert(), nil
}
