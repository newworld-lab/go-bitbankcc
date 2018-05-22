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

type transactionsResponse struct {
	baseResponse
	Data struct {
		Transactions []entity.Transaction `json:"transactions"`
	} `json:"data"`
}

type API interface {
	GetTicker(pair constant.TypePair) (*entity.Ticker, error)
	GetDepth(pair constant.TypePair) (*entity.Depth, error)
	GetTransactions(pair constant.TypePair, time *time.Time) (*entity.Transaction, error)
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

func (api *APIImpl) GetTransactions(pair constant.TypePair, t *time.Time, option *APIOption) ([]entity.Transaction, error) {
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

	tmp := new(struct {
		baseResponse
		Data struct {
			Transactions []struct {
				TransactionId int     `json:"transaction_id"`
				Side          string  `json:"side"`
				Price         float64 `json:"price,string"`
				Amount        float64 `json:"amount,string"`
				ExecutedAt    int64   `json:"executed_at"`
			}
		} `json:"data"`
	})
	json.Unmarshal(bytes, tmp)

	res := new(transactionsResponse)

	res.baseResponse = tmp.baseResponse
	for _, transacion := range tmp.Data.Transactions {
		res.Data.Transactions = append(res.Data.Transactions, entity.Transaction{
			TransactionId: transacion.TransactionId,
			Side:          transacion.Side,
			Price:         transacion.Price,
			Amount:        transacion.Amount,
			ExecutedAt:    time.Unix(transacion.ExecutedAt/1000, transacion.ExecutedAt%1000*1000000),
		})
	}

	err = tmp.parseError()
	if err != nil {
		return nil, err
	}

	return res.Data.Transactions, nil

}
