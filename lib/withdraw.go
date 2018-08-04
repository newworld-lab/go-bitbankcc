package lib

import (
	"encoding/json"
	"fmt"
	"net/http"

	entity "github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/pkg/errors"
)

type withdrawResponse struct {
	baseResponse
	Data struct {
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
			Uuid:    a.Account.Uuid,
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
	json.Unmarshal(bytes, res)

	err = res.parseError()
	if err != nil {
		return nil, err
	}

	return res.Data.Accounts.convert(), nil
}
