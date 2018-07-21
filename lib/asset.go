package lib

import (
	"encoding/json"
	"net/http"
	"strconv"

	constant "github.com/newworld-lab/go-bitbankcc/constant"
	entity "github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/pkg/errors"
)

type assetsResponse struct {
	baseResponse
	Data struct {
		Assets assets `json:"assets"`
	} `json:"data"`
}

type assets []asset

type asset struct {
	entity.Asset
	OnhandAmount string `json:"onhand_amount"`
	LockedAmount string `json:"locked_amount"`
	FreeAmount   string `json:"free_amount"`
}

func (as assets) convert() entity.Assets {
	assets := make(entity.Assets, 0)
	for _, a := range as {
		o, _ := strconv.ParseFloat(a.OnhandAmount, 64)
		l, _ := strconv.ParseFloat(a.LockedAmount, 64)
		f, _ := strconv.ParseFloat(a.FreeAmount, 64)
		assets = append(assets, entity.Asset{
			Asset:           a.Asset.Asset,
			AmountPrecision: a.AmountPrecision,
			OnhandAmount:    o,
			LockedAmount:    l,
			FreeAmount:      f,
			WithDrawalFee:   a.Asset.WithDrawalFee,
		})
	}
	return assets
}

func (api *APIImpl) GetAssets() (entity.Assets, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	if api.option == nil || api.option.ApiKey == nil || api.option.ApiSecret == nil {
		return nil, errors.New("ApiKey or ApiSecret is nil")
	}

	path := formatAssets
	header, err := api.createCertificationHeader(path)
	if err != nil {
		return nil, err
	}

	bytes, err := api.client.request(&clientOption{
		endpoint: constant.PrivateApiEndpoint,
		method:   http.MethodGet,
		path:     path,
		header:   header,
	})

	if err != nil {
		return nil, err
	}

	res := new(assetsResponse)
	json.Unmarshal(bytes, res)

	err = res.parseError()
	if err != nil {
		return nil, err
	}

	return res.Data.Assets.convert(), nil
}
