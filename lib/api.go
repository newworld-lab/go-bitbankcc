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
	GetAssets()
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

type tickerResponse struct {
	baseResponse
	Data struct {
		entity.Ticker
	} `json:"data"`
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

func (api *APIImpl) GetTransactions(pair constant.TypePair, t *time.Time, option *APIOption) (entity.Transactions, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	var path string
	if t == nil {
		path = fmt.Sprintf(formatTransactionsAll, pair)
	} else {
		fmt.Println(t.Format("20060102"))
		path = fmt.Sprintf(formatTransactions, pair, t.Format("20060102"))
	}
	bytes, err := api.client.request(&clientOption{
		method:  http.MethodGet,
		path:    path,
		timeout: option.GetTimeout(),
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
