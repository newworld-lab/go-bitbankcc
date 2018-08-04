package lib

import (
	"encoding/json"
	"net/http"
	"net/url"
	time "time"

	"github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/pkg/errors"
)

type orderResponse struct {
	baseResponse
	Data order `json:"data"`
}

type order struct {
	baseData
	entity.Order
	OrderedAt  int64  `json:"ordered_at"`
	ExecutedAt *int64 `json:"executed_at,omitempty"`
}

func (o *order) convert() *entity.Order {
	var (
		orderedAt  time.Time
		executedAt *time.Time
	)
	orderedAt = time.Unix(o.OrderedAt/1000, o.OrderedAt%1000*1000000)
	if o.ExecutedAt != nil {
		e := time.Unix(*o.ExecutedAt/1000, *o.ExecutedAt%1000*1000000)
		executedAt = &e
	}

	return &entity.Order{
		OrderID:         o.OrderID,
		Pair:            o.Pair,
		Type:            o.Type,
		StartAmount:     o.StartAmount,
		RemainingAmount: o.RemainingAmount,
		ExecutedAmount:  o.ExecutedAmount,
		Price:           o.Price,
		AveragePrice:    o.AveragePrice,
		OrderedAt:       orderedAt,
		ExecutedAt:      executedAt,
		Status:          o.Status,
	}
}

func (api *APIImpl) GetOrder(params entity.GetOrderParams) (*entity.Order, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	if api.option == nil || api.option.ApiKey == nil || api.option.ApiSecret == nil {
		return nil, errors.New("ApiKey or ApiSecret is nil")
	}

	url := &url.URL{}
	query := url.Query()
	query.Add("pair", string(params.Pair))
	query.Add("order_id", params.OrderID)
	query.Encode()
	path := formatOrder + "?" + query.Encode()
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

	res := new(orderResponse)
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}

	if res.Success != 1 {
		return nil, errors.Errorf("api error code=%d", res.Data.Code)
	}

	return res.Data.convert(), nil
}

func (api *APIImpl) PostOrder(params entity.PostOrderParams) (*entity.Order, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	if api.option == nil || api.option.ApiKey == nil || api.option.ApiSecret == nil {
		return nil, errors.New("ApiKey or ApiSecret is nil")
	}
	path := formatOrder

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

	res := new(orderResponse)
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}

	if res.Success != 1 {
		return nil, errors.Errorf("api error code=%d", res.Data.Code)
	}

	return res.Data.convert(), nil
}
