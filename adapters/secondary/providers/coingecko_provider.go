package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/moondolphin/crypto-api/domain"
)

type CoinGeckoProvider struct {
	BaseURL string
	Client  *http.Client
}

func NewCoinGeckoProvider() *CoinGeckoProvider {
	return &CoinGeckoProvider{
		BaseURL: "https://api.coingecko.com/api/v3",
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *CoinGeckoProvider) Name() string { return "coingecko" }

func (p *CoinGeckoProvider) GetCurrentPrice(ctx context.Context, coin domain.Coin, currency string) (domain.PriceQuote, error) {
	id := strings.TrimSpace(coin.CoinGeckoID)
	if id == "" {
		return domain.PriceQuote{}, fmt.Errorf("coingecko_id missing for symbol %s", coin.Symbol)
	}

	vs := strings.ToLower(strings.TrimSpace(currency))
	if vs == "" {
		return domain.PriceQuote{}, fmt.Errorf("currency required")
	}

	endpoint, err := url.Parse(p.BaseURL + "/simple/price")
	if err != nil {
		return domain.PriceQuote{}, err
	}

	q := endpoint.Query()
	q.Set("ids", id)
	q.Set("vs_currencies", vs)
	endpoint.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return domain.PriceQuote{}, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return domain.PriceQuote{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return domain.PriceQuote{}, fmt.Errorf("coingecko status %d", resp.StatusCode)
	}

	// Response example: { "bitcoin": { "usd": 88338 } }
	var raw map[string]map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return domain.PriceQuote{}, err
	}

	m, ok := raw[id]
	if !ok {
		return domain.PriceQuote{}, fmt.Errorf("coingecko missing id %s", id)
	}

	v, ok := m[vs]
	if !ok {
		return domain.PriceQuote{}, fmt.Errorf("coingecko missing currency %s", vs)
	}

	priceStr, err := toDecimalString(v)
	if err != nil {
		return domain.PriceQuote{}, err
	}

	return domain.PriceQuote{
		Symbol:    coin.Symbol,
		Currency:  strings.ToUpper(currency),
		Price:     priceStr,
		Provider:  p.Name(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func toDecimalString(v any) (string, error) {
	switch t := v.(type) {
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case json.Number:
		return t.String(), nil
	case string:
		return t, nil
	default:
		return "", fmt.Errorf("unexpected price type %T", v)
	}
}
