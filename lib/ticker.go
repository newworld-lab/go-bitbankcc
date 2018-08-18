package lib

import (
	"encoding/json"
	"fmt"
	"net/http"

	entity "github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/pkg/errors"
)

const (
	formatTicker = "/%s/ticker"
)

type tickerResponse struct {
	baseResponse
	Data struct {
		baseData
		entity.Ticker
	} `json:"data"`
}

func (api *APIImpl) GetTicker(pair entity.TypePair) (*entity.Ticker, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	bytes, err := api.client.request(&clientOption{
		endpoint: publicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatTicker, pair),
	})
	if err != nil {
		return nil, err
	}

	res := new(tickerResponse)
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Success != 1 {
		return nil, errors.Errorf("api error code=%d", res.Data.Code)
	}

	return &res.Data.Ticker, nil
}
