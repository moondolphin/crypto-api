package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/moondolphin/crypto-api/domain"
)

var ErrBinanceAPI = errors.New("binance_api_error")

type BinanceProvider struct {
	BaseURL string
	Client  *http.Client
}

func NewBinanceProvider() *BinanceProvider {
	return &BinanceProvider{
		BaseURL: "https://api.binance.com",
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (p *BinanceProvider) Name() string {
	return "binance"
}

func (p *BinanceProvider) GetCurrentPrice(
	ctx context.Context,
	coin domain.Coin,
	currency string,
) (domain.PriceQuote, error) {

	if coin.BinanceSymbol == "" {
		return domain.PriceQuote{}, ErrBinanceAPI
	}

	// Binance usa símbolo de par, ej: BTCUSDT
	url := fmt.Sprintf(
		"%s/api/v3/ticker/price?symbol=%s",
		p.BaseURL,
		coin.BinanceSymbol,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.PriceQuote{}, err
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return domain.PriceQuote{}, ErrBinanceAPI
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.PriceQuote{}, ErrBinanceAPI
	}

	var r struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return domain.PriceQuote{}, err
	}

	return domain.PriceQuote{
		Symbol:   coin.Symbol,
		Currency: currency,
		Price:    r.Price,
		Provider: p.Name(),
	}, nil
}
