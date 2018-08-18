package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	entity "github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/pkg/errors"
)

const (
	formatDepth = "/%s/depth"
)

type depth struct {
	entity.Depth
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

type depthResponse struct {
	baseResponse
	Data struct {
		depth
	} `json:"data"`
}

func (d *depth) convert() entity.Depth {
	asks, bids := make([][]float64, 0), make([][]float64, 0)
	for _, ss := range d.Asks {
		fs := make([]float64, 0)
		for _, s := range ss {
			f, _ := strconv.ParseFloat(s, 64)
			fs = append(fs, f)
		}
		asks = append(asks, fs)
	}
	for _, ss := range d.Bids {
		fs := make([]float64, 0)
		for _, s := range ss {
			f, _ := strconv.ParseFloat(s, 64)
			fs = append(fs, f)
		}
		bids = append(bids, fs)
	}
	return entity.Depth{
		Asks: asks,
		Bids: bids,
	}
}

func (api *APIImpl) GetDepth(pair entity.TypePair) (*entity.Depth, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	bytes, err := api.client.request(&clientOption{
		endpoint: publicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatDepth, pair),
	})
	if err != nil {
		return nil, err
	}

	res := new(depthResponse)
	json.Unmarshal(bytes, res)
	err = res.parseError()
	if err != nil {
		return nil, err
	}

	depth := res.Data.convert()
	return &depth, nil
}
