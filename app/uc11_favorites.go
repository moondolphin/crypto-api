package app

type FavoriteCoinOutput struct {
	ID            int64  `json:"id"`
	Symbol        string `json:"symbol"`
	Enabled       bool   `json:"enabled"`
	CoinGeckoID   string `json:"coingecko_id"`
	BinanceSymbol string `json:"binance_symbol"`
}
