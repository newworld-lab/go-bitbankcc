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

const (
	publicApiEndpoint  = "https://public.bitbank.cc"
	privateApiEndpoint = "https://api.bitbank.cc"
)

const (
	formatAccessSignature = "%d%s%s"
)

type baseData struct {
	Code int `json:"code"`
}

type baseResponse struct {
	Success int      `json:"success"`
	Data    baseData `json:"data"`
}

type API interface {
	GetTicker(pair entity.TypePair) (*entity.Ticker, error)
	GetDepth(pair entity.TypePair) (*entity.Depth, error)
	GetTransactions(pair entity.TypePair, time *time.Time) (*entity.Transaction, error)
	GetCandlestick(pair entity.TypePair, candle TypeCandle, time time.Time) (entity.Candlestick, error)
	GetAssets() (entity.Assets, error)
	GetWithdraw(asset entity.TypeAsset) (*entity.Accounts, error)
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

func (api *APIImpl) createCertificationHeader(str string) (http.Header, error) {
	accessNonce := api.createAccessNonce()
	mac := hmac.New(sha256.New, []byte(*api.option.ApiSecret))
	_, err := mac.Write([]byte(fmt.Sprintf(formatAccessSignature, accessNonce, str, "")))
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
