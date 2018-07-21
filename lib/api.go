package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/newworld-lab/go-bitbankcc/entity"

	"github.com/pkg/errors"
)

type TypePair string

const (
	PairBtcJpy  TypePair = "btc_jpy"
	PairXrpJpy  TypePair = "xrp_jpy"
	PairLtcBtc  TypePair = "ltc_btc"
	PairEthBtc  TypePair = "eth_btc"
	PairMonaJpy TypePair = "mona_jpy"
	PairMonaBtc TypePair = "mona_btc"
	PairBccJpy  TypePair = "bcc_jpy"
	PairBccBtc  TypePair = "bcc_btc"
)

const (
	publicApiEndpoint  = "https://public.bitbank.cc"
	privateApiEndpoint = "https://api.bitbank.cc"
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
	GetTicker(pair TypePair) (*entity.Ticker, error)
	GetDepth(pair TypePair) (*entity.Depth, error)
	GetTransactions(pair TypePair, time *time.Time) (*entity.Transaction, error)
	GetCandlestick(pair TypePair, candle TypeCandle, time time.Time) (entity.Candlestick, error)
	GetAssets() (entity.Assets, error)
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

func (api *APIImpl) createCertificationHeader(path string) (http.Header, error) {
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
	return header, nil
}
