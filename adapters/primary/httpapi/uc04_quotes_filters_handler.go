package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/app"
)

type GetQuoteFiltersHandler struct {
	UC app.GetQuoteFiltersUseCase
}

// @Summary Obtener filtros disponibles (faceted filters) para cotizaciones
// @Description Devuelve los valores disponibles para combos (symbol/provider/currency) y rangos, recalculados según los filtros aplicados.
// @Tags Quotes
// @Param symbol query string false "Símbolo (BTC, ETH...)"
// @Param provider query string false "Proveedor (binance, coingecko...)"
// @Param currency query string false "Moneda (USD, USDT...)"
// @Param min_price query number false "Precio mínimo"
// @Param max_price query number false "Precio máximo"
// @Param from query string false "Desde. Formatos: 'YYYY-MM-DD' o 'YYYY-MM-DDTHH:MM:SSZ'"
// @Param to   query string false "Hasta. Formatos: 'YYYY-MM-DD' o 'YYYY-MM-DDTHH:MM:SSZ'"
// @Success 200 {object} app.GetQuoteFiltersOutput
// @Failure 400 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/v1/quotes/filters [get]
func (h GetQuoteFiltersHandler) Handle(c *gin.Context) {
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

	// from / to (RFC3339 o YYYY-MM-DD) - reutiliza parseTimeFlexible del search handler
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

	out, err := h.UC.Execute(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusOK, out)
}
