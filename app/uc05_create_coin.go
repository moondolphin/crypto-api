package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/moondolphin/crypto-api/domain"
)

var (
	ErrInvalidCoinInput  = errors.New("invalid_coin_input")
	ErrCoinNotResolvable = errors.New("coin_not_resolvable")
)

type CreateCoinInput struct {
	Symbol        string `json:"symbol"`
	Enabled       *bool  `json:"enabled,omitempty"`
	CoinGeckoID   string `json:"coingecko_id,omitempty"`
	BinanceSymbol string `json:"binance_symbol,omitempty"`
}

type CreateCoinOutput struct {
	ID            int64  `json:"id"`
	Symbol        string `json:"symbol"`
	Enabled       bool   `json:"enabled"`
	CoinGeckoID   string `json:"coingecko_id"`
	BinanceSymbol string `json:"binance_symbol"`
}


type ProviderGetter interface {
	Get(name string) (domain.PriceProvider, bool)
}

type CreateCoinUseCase struct {
	CoinRepo  domain.CoinRepository
	Providers ProviderGetter

	// Preferido para Binance (por defecto es USDT)
	BinanceQuoteCurrency string

	// HTTP client inyectable
	HTTPClient *http.Client
}

func (uc CreateCoinUseCase) Execute(ctx context.Context, in CreateCoinInput) (CreateCoinOutput, error) {
	symbol := strings.ToUpper(strings.TrimSpace(in.Symbol))
	if symbol == "" {
		return CreateCoinOutput{}, ErrInvalidCoinInput
	}

	enabled := true
	if in.Enabled != nil {
		enabled = *in.Enabled
	}

	client := uc.httpClient()

	// Inputs(Swagger "string" -> vacío)
	inCG := sanitizeOptionalString(in.CoinGeckoID)
	inBN := sanitizeOptionalString(in.BinanceSymbol)

	
	// Si son inválidos: se ignoran (no pisan) y luego se auto-resuelven.
	if inBN != "" {
		if err := validateBinanceSymbol(ctx, client, inBN); err != nil {
			inBN = ""
		}
	}
	if inCG != "" {
		if err := validateCoinGeckoID(ctx, client, inCG); err != nil {
			inCG = ""
		}
	}

	// MERGE: traemos existente para no pisar con vacío
	existing, err := uc.CoinRepo.GetBySymbol(ctx, symbol)
	if err != nil {
		return CreateCoinOutput{}, err
	}

	merged := domain.Coin{
		Symbol:  symbol,
		Enabled: enabled,
	}

	if existing != nil {
		merged.ID = existing.ID
		merged.CoinGeckoID = sanitizeOptionalString(existing.CoinGeckoID)
		merged.BinanceSymbol = sanitizeOptionalString(existing.BinanceSymbol)
	}

	// Aplicar overrides SOLO si quedaron valores válidos desde el input
	if inCG != "" {
		merged.CoinGeckoID = inCG
	}
	if inBN != "" {
		merged.BinanceSymbol = inBN
	}

	// Si aún no tenemos IDs, auto-resolve (symbol-only o inputs inválidos)
	if merged.CoinGeckoID == "" && merged.BinanceSymbol == "" {
		if err := uc.autoResolve(ctx, client, &merged); err != nil {
			return CreateCoinOutput{}, ErrCoinNotResolvable
		}
	}

	out, err := uc.CoinRepo.Upsert(ctx, merged)
	if err != nil {
		return CreateCoinOutput{}, err
	}

	return CreateCoinOutput{
		ID:            out.ID,
		Symbol:        out.Symbol,
		Enabled:       out.Enabled,
		CoinGeckoID:   out.CoinGeckoID,
		BinanceSymbol: out.BinanceSymbol,
	}, nil
}

func (uc CreateCoinUseCase) httpClient() *http.Client {
	if uc.HTTPClient != nil {
		return uc.HTTPClient
	}
	return &http.Client{Timeout: 8 * time.Second}
}

func sanitizeOptionalString(s string) string {
	v := strings.TrimSpace(s)
	if v == "" {
		return ""
	}
	if strings.EqualFold(v, "string") {
		return ""
	}
	return v
}

func (uc CreateCoinUseCase) autoResolve(ctx context.Context, client *http.Client, coin *domain.Coin) error {
	quote := strings.TrimSpace(uc.BinanceQuoteCurrency)
	if quote == "" {
		quote = "USDT"
	}

	// 1) Binance: probar SYMBOL+USDT si registry tiene binance
	if coin.BinanceSymbol == "" {
		if _, ok := uc.Providers.Get("binance"); ok {
			pair := strings.ToUpper(strings.TrimSpace(coin.Symbol) + quote)
			if bn, err := resolveBinancePair(ctx, client, pair); err == nil {
				coin.BinanceSymbol = bn
			}
		}
	}

	// 2) CoinGecko: buscar ID por symbol si registry tiene coingecko
	if coin.CoinGeckoID == "" {
		if _, ok := uc.Providers.Get("coingecko"); ok {
			if id, err := resolveCoinGeckoID(ctx, client, coin.Symbol); err == nil {
				coin.CoinGeckoID = id
			}
		}
	}

	if coin.CoinGeckoID == "" && coin.BinanceSymbol == "" {
		return fmt.Errorf("could not resolve coin ids for %s", coin.Symbol)
	}
	return nil
}

// =====================
// Validaciones
// =====================

func validateBinanceSymbol(ctx context.Context, client *http.Client, pair string) error {
	_, err := resolveBinancePair(ctx, client, strings.ToUpper(strings.TrimSpace(pair)))
	return err
}

func validateCoinGeckoID(ctx context.Context, client *http.Client, id string) error {
	// validate usando /simple/price con vs=usd
	baseURL := "https://api.coingecko.com/api/v3"
	endpoint, err := url.Parse(baseURL + "/simple/price")
	if err != nil {
		return err
	}
	q := endpoint.Query()
	q.Set("ids", strings.TrimSpace(id))
	q.Set("vs_currencies", "usd")
	endpoint.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("coingecko status %d", resp.StatusCode)
	}

	var raw map[string]map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return err
	}
	if _, ok := raw[strings.TrimSpace(id)]; !ok {
		return fmt.Errorf("coingecko invalid id %s", id)
	}
	return nil
}

// =====================
// Resoluciones
// =====================

func resolveBinancePair(ctx context.Context, client *http.Client, pair string) (string, error) {
	baseURL := "https://api.binance.com"
	u := fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", baseURL, url.QueryEscape(pair))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("binance status %d", resp.StatusCode)
	}

	var r struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	if strings.TrimSpace(r.Symbol) == "" {
		return "", fmt.Errorf("binance empty symbol")
	}
	return strings.ToUpper(strings.TrimSpace(r.Symbol)), nil
}

func resolveCoinGeckoID(ctx context.Context, client *http.Client, symbol string) (string, error) {
	baseURL := "https://api.coingecko.com/api/v3"
	endpoint, err := url.Parse(baseURL + "/search")
	if err != nil {
		return "", err
	}
	q := endpoint.Query()
	q.Set("query", symbol)
	endpoint.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("coingecko status %d", resp.StatusCode)
	}

	var raw struct {
		Coins []struct {
			ID     string `json:"id"`
			Symbol string `json:"symbol"`
			Name   string `json:"name"`
		} `json:"coins"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return "", err
	}

	sym := strings.ToLower(strings.TrimSpace(symbol))
	for _, c := range raw.Coins {
		if strings.ToLower(strings.TrimSpace(c.Symbol)) == sym && strings.TrimSpace(c.ID) != "" {
			return strings.TrimSpace(c.ID), nil
		}
	}

	if len(raw.Coins) > 0 && strings.TrimSpace(raw.Coins[0].ID) != "" {
		return strings.TrimSpace(raw.Coins[0].ID), nil
	}

	return "", fmt.Errorf("coingecko id not found for %s", symbol)
}
