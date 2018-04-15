package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KoteiIto/go-bitbankcc/constant"
	"github.com/KoteiIto/go-bitbankcc/entity"

	"github.com/pkg/errors"
)

const (
	formatTicker = "/%s/ticker"
)

type baseResponse struct {
	Success int `json:"success"`
	Data    struct {
		Code int `json:"code"`
	} `json:"data"`
}

func (res *baseResponse) parseError() error {
	if res == nil {
		return errors.New("res is nil")
	}

	if res.Success != 1 {
		return errors.Errorf("api error code=%d", res.Data.Code)
	}

	return nil
}

type tickerResponse struct {
	baseResponse
	Data struct {
		entity.Ticker
	} `json:"data"`
}

type API interface {
	GetTicker(pair constant.TypePair) (*entity.Ticker, error)
}

type APIImpl struct {
	client *clientImpl
}

type APIOption struct {
	timeout time.Duration
}

func (option *APIOption) GetTimeout() time.Duration {
	if option == nil {
		return time.Duration(0)
	}
	return option.timeout
}

func (api *APIImpl) GetTicker(pair constant.TypePair, option *APIOption) (*entity.Ticker, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	bytes, err := api.client.request(&clientOption{
		method:  http.MethodGet,
		path:    fmt.Sprintf(formatTicker, pair),
		timeout: option.GetTimeout(),
	})
	if err != nil {
		return nil, err
	}

	res := new(tickerResponse)
	json.Unmarshal(bytes, res)

	err = res.parseError()
	if err != nil {
		return nil, err
	}

	return &res.Data.Ticker, nil
}
