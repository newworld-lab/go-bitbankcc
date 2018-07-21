package lib

import (
	"encoding/json"
	"fmt"
	"net/http"

	entity "github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/pkg/errors"
)

type tickerResponse struct {
	baseResponse
	Data struct {
		entity.Ticker
	} `json:"data"`
}

func (api *APIImpl) GetTicker(pair TypePair) (*entity.Ticker, error) {
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
	json.Unmarshal(bytes, res)

	err = res.parseError()
	if err != nil {
		return nil, err
	}

	return &res.Data.Ticker, nil
}
