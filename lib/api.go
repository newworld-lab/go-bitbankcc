package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/newworld-lab/go-bitbankcc/constant"
	"github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/newworld-lab/go-bitbankcc/util"

	"github.com/pkg/errors"
)

const (
	formatTicker = "/%s/ticker"
	formatDepth  = "/%s/depth"
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

type depthResponse struct {
	baseResponse
	Data struct {
		entity.Depth
	} `json:"data"`
}

type API interface {
	GetTicker(pair constant.TypePair) (*entity.Ticker, error)
	GetDepth(pair constant.TypePair) (*entity.Depth, error)
}

type APIImpl struct {
	client client
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

// GetTicker 通貨TypeからTicker取得
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

// GetDepth 通貨TypeからDepth取得
func (api *APIImpl) GetDepth(pair constant.TypePair, option *APIOption) (*entity.Depth, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	bytes, err := api.client.request(&clientOption{
		method:  http.MethodGet,
		path:    fmt.Sprintf(formatDepth, pair),
		timeout: option.GetTimeout(),
	})
	if err != nil {
		return nil, err
	}

	// 配列をFloatにキャストするための仮struct
	tmp := new(struct {
		baseResponse
		Data struct {
			Asks []util.Strings `json:"asks"`
			Bids []util.Strings `json:"bids"`
		} `json:"data"`
	})
	json.Unmarshal(bytes, tmp)

	res := new(depthResponse)
	// resの中身をキャストした内容で書き換え
	res.baseResponse = tmp.baseResponse
	for _, strAsks := range tmp.Data.Asks {
		res.Data.Asks = append(res.Data.Asks, strAsks.ToFloat64())
	}
	for _, strBids := range tmp.Data.Bids {
		res.Data.Bids = append(res.Data.Bids, strBids.ToFloat64())
	}

	err = res.parseError()
	if err != nil {
		return nil, err
	}

	return &res.Data.Depth, nil
}
