package domain

type Coin struct {
	ID            int64
	Symbol        string
	Enabled       bool
	CoinGeckoID   string // ej: "bitcoin"
	BinanceSymbol string // ej: "BTCUSDT"
}
