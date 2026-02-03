package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/app"
)

type SearchQuotesHandler struct {
	UC app.SearchQuotesUseCase
}

// @Summary Listar cotizaciones (histórico) con filtros, paginado y summary
// @Description Devuelve cotizaciones persistidas en BD. Permite filtros combinables y paginado (máx 10 páginas).
// @Tags Quotes
// @Param symbol query string false "Símbolo (BTC, ETH , ADA, APT, ATOM, AAVE...)"
// @Param provider query string false "Proveedor (binance, coingecko...)"
// @Param currency query string false "Moneda (USD, USDT...)"
// @Param min_price query number false "Precio mínimo"
// @Param max_price query number false "Precio máximo"
// @Param from query string false "Desde. Formatos: 'YYYY-MM-DD' o 'YYYY-MM-DDTHH:MM:SSZ'. Ej: 2026-01-30 o 2026-01-30T10:00:00Z"
// @Param to   query string false "Hasta. Formatos: 'YYYY-MM-DD' o 'YYYY-MM-DDTHH:MM:SSZ'. Ej: 2026-01-30 o 2026-01-30T23:59:59Z"
// @Param page query int false "Página (1..10)"
// @Param page_size query int false "Tamaño (1..100)"
// @Success 200 {object} app.SearchQuotesOutput
// @Failure 400 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/v1/quotes [get]
func (h SearchQuotesHandler) Handle(c *gin.Context) {
	in := app.SearchQuotesInput{
		Symbol:   c.Query("symbol"),
		Provider: c.Query("provider"),
		Currency: c.Query("currency"),
	}

	// min_price / max_price
	if v := strings.TrimSpace(c.Query("min_price")); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_min_price"})
			return
		}
		in.MinPrice = &f
	}
	if v := strings.TrimSpace(c.Query("max_price")); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_max_price"})
			return
		}
		in.MaxPrice = &f
	}

	// from / to (RFC3339 o YYYY-MM-DD)
	if v := strings.TrimSpace(c.Query("from")); v != "" {
		tm, err := parseTimeFlexible(v, false)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_from"})
			return
		}
		in.From = &tm
	}
	if v := strings.TrimSpace(c.Query("to")); v != "" {
		tm, err := parseTimeFlexible(v, true)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_to"})
			return
		}
		in.To = &tm
	}

	// page / page_size
	if v := strings.TrimSpace(c.Query("page")); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_page"})
			return
		}
		in.Page = n
	}
	if v := strings.TrimSpace(c.Query("page_size")); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_page_size"})
			return
		}
		in.PageSize = n
	}

	out, err := h.UC.Execute(c.Request.Context(), in)
	if err != nil {
		switch err {
		case app.ErrInvalidFilters:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
		}
		return
	}

	c.JSON(http.StatusOK, out)
}

func parseTimeFlexible(s string, isTo bool) (time.Time, error) {
	// 1) RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	// 2) YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", s); err == nil {
		t = t.UTC()
		if isTo {
			// fin del día (incluyente) con precisión microsegundos (DATETIME(6))
			t = t.Add(24*time.Hour - time.Microsecond)
		}
		return t, nil
	}

	return time.Time{}, errors.New("invalid_time_format")
}
