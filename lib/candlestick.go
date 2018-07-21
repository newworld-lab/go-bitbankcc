package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	time "time"

	constant "github.com/newworld-lab/go-bitbankcc/constant"
	entity "github.com/newworld-lab/go-bitbankcc/entity"
	"github.com/pkg/errors"
)

type TypeCandle string

const (
	OneMinute      TypeCandle = "1min"
	FiveMinutes    TypeCandle = "5min"
	FifteenMinutes TypeCandle = "15min"
	ThirtyMinutes  TypeCandle = "30min"
	OneHour        TypeCandle = "1hour"
	FourHours      TypeCandle = "4hour"
	EightHours     TypeCandle = "8hour"
	TwelveHours    TypeCandle = "12hour"
	OneDay         TypeCandle = "1day"
	OneWeek        TypeCandle = "1week"
)

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

func (api *APIImpl) GetCandlestick(pair TypePair, candle TypeCandle, t time.Time) (entity.Candlestick, error) {
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
