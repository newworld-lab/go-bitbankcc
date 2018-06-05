package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/newworld-lab/go-bitbankcc/constant"
	"github.com/newworld-lab/go-bitbankcc/entity"

	"github.com/pkg/errors"
)

const (
	formatTicker          = "/%s/ticker"
	formatDepth           = "/%s/depth"
	formatTransactionsAll = "/%s/transactions"
	formatTransactions    = "/%s/transactions/%s"
	formatCandlestick     = "/%s/candlestick/%s/%s"
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

type API interface {
	GetTicker(pair constant.TypePair) (*entity.Ticker, error)
	GetDepth(pair constant.TypePair) (*entity.Depth, error)
	GetTransactions(pair constant.TypePair, time *time.Time) (*entity.Transaction, error)
	GetCandlestick(pair constant.TypePair, candle constant.TypeCandle, time time.Time) (entity.Candlestick, error)
	GetAssets()
}

type APIImpl struct {
	client client
	option *APIOption
}

func NewApi(option *APIOption) *APIImpl {
	timeout := 5000 * time.Millisecond
	if option != nil && option.timeout != 0 {
		timeout = option.timeout
	}
	return &APIImpl{
		client: &clientImpl{
			httpClient: http.Client{
				Timeout: timeout,
			},
		},
		option: option,
	}
}

type APIOption struct {
	timeout time.Duration
}

type tickerResponse struct {
	baseResponse
	Data struct {
		entity.Ticker
	} `json:"data"`
}

// GetTicker 通貨TypeからTicker取得
func (api *APIImpl) GetTicker(pair constant.TypePair) (*entity.Ticker, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	bytes, err := api.client.request(&clientOption{
		endpoint: constant.PublicApiEndpoint,
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

// GetDepth 通貨TypeからDepth取得
func (api *APIImpl) GetDepth(pair constant.TypePair) (*entity.Depth, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	bytes, err := api.client.request(&clientOption{
		endpoint: constant.PublicApiEndpoint,
		method:   http.MethodGet,
		path:     fmt.Sprintf(formatDepth, pair),
	})
	if err != nil {
		return nil, err
	}

	// 配列をFloatにキャストするための仮struct
	res := new(depthResponse)
	json.Unmarshal(bytes, res)
	err = res.parseError()
	if err != nil {
		return nil, err
	}

	depth := res.Data.convert()
	return &depth, nil
}

type transactionsResponse struct {
	baseResponse
	Data struct {
		Transactions transactions `json:"transactions"`
	} `json:"data"`
}

type transaction struct {
	entity.Transaction
	ExecutedAt int64 `json:"executed_at"`
}

type transactions []transaction

func (ts transactions) convert() entity.Transactions {
	transactions := make(entity.Transactions, 0)
	for _, t := range ts {
		transactions = append(transactions, entity.Transaction{
			TransactionId: t.TransactionId,
			Side:          t.Side,
			Price:         t.Price,
			Amount:        t.Amount,
			ExecutedAt:    time.Unix(t.ExecutedAt/1000, t.ExecutedAt%1000*1000000),
		})
	}
	return transactions
}

func (api *APIImpl) GetTransactions(pair constant.TypePair, t *time.Time) (entity.Transactions, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	var path string
	if t == nil {
		path = fmt.Sprintf(formatTransactionsAll, pair)
	} else {
		path = fmt.Sprintf(formatTransactions, pair, t.Format("20060102"))
	}
	bytes, err := api.client.request(&clientOption{
		endpoint: constant.PublicApiEndpoint,
		method:   http.MethodGet,
		path:     path,
	})

	if err != nil {
		return nil, err
	}

	res := new(transactionsResponse)
	json.Unmarshal(bytes, res)

	err = res.parseError()
	if err != nil {
		return nil, err
	}

	return res.Data.Transactions.convert(), nil
}

type candlestickResponse struct {
	baseResponse
	Data struct {
		Candlestick candlestick `json:"candlestick"`
	} `json:"data"`
}

type candlestick []struct {
	entity.CandlestickItem
	Ohlcv [][]interface{} `json:"ohlcv"` // [][]だったのを[]にしてみる
}

func (c candlestick) convert() entity.Candlestick {

	candlestick := make(entity.Candlestick, 0)
	for _, i := range c {
		candlestickItem := entity.CandlestickItem{
			Type: i.Type,
		}
		for _, item := range i.Ohlcv {
			ohlcvItem := entity.OhlcvItem{}
			ohlcvItem.Open, _ = strconv.ParseFloat(item[0].(string), 64)
			ohlcvItem.High, _ = strconv.ParseFloat(item[1].(string), 64)
			ohlcvItem.Low, _ = strconv.ParseFloat(item[2].(string), 64)
			ohlcvItem.Close, _ = strconv.ParseFloat(item[3].(string), 64)
			ohlcvItem.Volume, _ = strconv.ParseFloat(item[4].(string), 64)
			intDate := int64(item[5].(float64))
			ohlcvItem.Date = time.Unix(intDate/1000, intDate%1000*1000000)
			candlestickItem.Ohlcv = append(candlestickItem.Ohlcv, ohlcvItem)
		}
		candlestick = append(candlestick, candlestickItem)
	}
	return candlestick
}

func (api *APIImpl) GetCandlestick(pair constant.TypePair, candle constant.TypeCandle, t time.Time) (entity.Candlestick, error) {

	if api == nil {
		return nil, errors.New("api is nil")
	}

	var path string

	path = fmt.Sprintf(formatCandlestick, pair, candle, t.Format("20060102"))
	bytes, err := api.client.request(&clientOption{
		endpoint: constant.PublicApiEndpoint,
		method:   http.MethodGet,
		path:     path,
	})

	if err != nil {
		return nil, err
	}

	res := new(candlestickResponse)
	json.Unmarshal(bytes, res)

	err = res.parseError()
	if err != nil {
		return nil, err
	}

	return res.Data.Candlestick.convert(), nil
}
