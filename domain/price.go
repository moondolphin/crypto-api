package domain

type PriceQuote struct {
	Symbol    string
	Currency  string
	Price     string
	Provider  string
	Timestamp string // RFC3339
}
