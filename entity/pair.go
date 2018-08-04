package entity

type TypePair string
type TypeAsset string

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
	AssetBtc  TypeAsset = "btc"
	AssetXrp  TypeAsset = "xrp"
	AssetLtc  TypeAsset = "ltc"
	AssetEth  TypeAsset = "eth"
	AssetMona TypeAsset = "mona"
	AssetBcc  TypeAsset = "bcc"
)
