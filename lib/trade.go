package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	time "time"

	entity "github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/pkg/errors"
)

const (
	formatTrades = "/v1/user/spot/trade_history"
)

type tradesResponse struct {
	baseResponse
	Data struct {
		baseData
		Trades trades `json:"trades"`
	} `json:"data"`
}

type trades []trade

type trade struct {
	entity.Trade
	ExecutedAt int64 `json:"executed_at"`
}

func (tr trades) convert() entity.Trades {
	trades := make(entity.Trades, 0)
	for _, t := range tr {
		trades = append(trades, entity.Trade{
			TradeId:        t.Trade.TradeId,
			Pair:           t.Trade.Pair,
			OrderId:        t.Trade.OrderId,
			Side:           t.Trade.Side,
			Type:           t.Trade.Type,
			Amount:         t.Trade.Amount,
			Price:          t.Trade.Price,
			MakerTaker:     t.Trade.MakerTaker,
			FeeAmountBase:  t.Trade.FeeAmountBase,
			FeeAmountQuote: t.Trade.FeeAmountQuote,
			ExecutedAt:     time.Unix(int64(t.ExecutedAt)/1000, int64(t.ExecutedAt)%1000*1000000),
		})
	}
	return trades
}

func (api *APIImpl) GetTrades(params entity.TradeParams) (entity.Trades, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	if api.option == nil || api.option.ApiKey == nil || api.option.ApiSecret == nil {
		return nil, errors.New("ApiKey or ApiSecret is nil")
	}

	path := formatTrades

	url := &url.URL{}
	query := url.Query()

	if params.Pair != "" {
		query.Add("pair", fmt.Sprint(params.Pair))
	}

	if params.Count != 0 {
		query.Add("count", fmt.Sprint(params.Count))
	}

	if params.OrderId != 0 {
		query.Add("order_id", fmt.Sprint(params.OrderId))
	}

	if params.Since != nil {
		query.Add("since", fmt.Sprint(params.Since.Unix()*1000))
	}

	if params.End != nil {
		query.Add("end", fmt.Sprint(params.Since.Unix()*1000))
	}

	if params.Order != "" {
		query.Add("order", fmt.Sprint(params.Order))
	}

	path = path + "?" + query.Encode()

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

	res := new(tradesResponse)
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return nil, err
	}

	if res.Success != 1 {
		return nil, errors.Errorf("api error code=%d", res.Data.Code)
	}

	return res.Data.Trades.convert(), nil
}
