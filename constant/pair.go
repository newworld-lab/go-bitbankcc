package constant

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
