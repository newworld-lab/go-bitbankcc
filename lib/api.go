package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	formatAssets          = "/v1/user/assets"
	formatAccessSignature = "%d%s%s"
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
	client            client
	option            *APIOption
	createAccessNonce func() int
}

func createAccessNonce() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}

func NewApi(option *APIOption) *APIImpl {
	timeout := 5000 * time.Millisecond
	if option != nil && option.Timeout != 0 {
		timeout = option.Timeout
	}

	client := &clientImpl{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}

	return &APIImpl{
		client:            client,
		option:            option,
		createAccessNonce: createAccessNonce,
	}
}

type APIOption struct {
	Timeout   time.Duration
	ApiKey    *string
	ApiSecret *string
}

type tickerResponse struct {
	baseResponse
	Data struct {
		entity.Ticker
	} `json:"data"`
}

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
	Ohlcv [][]interface{} `json:"ohlcv"`
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

type assetsResponse struct {
	baseResponse
	Data struct {
		Assets assets `json:"assets"`
	} `json:"data"`
}

type assets []asset

type asset struct {
	entity.Asset
	OnhandAmount string `json:"onhand_amount"`
	LockedAmount string `json:"locked_amount"`
	FreeAmount   string `json:"free_amount"`
}

func (as assets) convert() entity.Assets {
	assets := make(entity.Assets, 0)
	for _, a := range as {
		o, _ := strconv.ParseFloat(a.OnhandAmount, 64)
		l, _ := strconv.ParseFloat(a.LockedAmount, 64)
		f, _ := strconv.ParseFloat(a.FreeAmount, 64)
		assets = append(assets, entity.Asset{
			Asset:           a.Asset.Asset,
			AmountPrecision: a.AmountPrecision,
			OnhandAmount:    o,
			LockedAmount:    l,
			FreeAmount:      f,
			WithDrawalFee:   a.Asset.WithDrawalFee,
		})
	}
	return assets
}

func (api *APIImpl) GetAssets() (entity.Assets, error) {
	if api == nil {
		return nil, errors.New("api is nil")
	}

	if api.option == nil || api.option.ApiKey == nil || api.option.ApiSecret == nil {
		return nil, errors.New("ApiKey or ApiSecret is nil")
	}

	path := formatAssets
	accessNonce := api.createAccessNonce()
	mac := hmac.New(sha256.New, []byte(*api.option.ApiSecret))
	_, err := mac.Write([]byte(fmt.Sprintf(formatAccessSignature, accessNonce, path, "")))
	if err != nil {
		return nil, errors.Cause(err)
	}
	accessSignature := hex.EncodeToString(mac.Sum(nil))
	header := make(http.Header)
	header["ACCESS-KEY"] = []string{*api.option.ApiKey}
	header["ACCESS-NONCE"] = []string{strconv.Itoa(accessNonce)}
	header["ACCESS-SIGNATURE"] = []string{accessSignature}

	bytes, err := api.client.request(&clientOption{
		endpoint: constant.PrivateApiEndpoint,
		method:   http.MethodGet,
		path:     path,
		header:   header,
	})

	if err != nil {
		return nil, err
	}

	res := new(assetsResponse)
	json.Unmarshal(bytes, res)

	err = res.parseError()
	if err != nil {
		return nil, err
	}

	return res.Data.Assets.convert(), nil
}
