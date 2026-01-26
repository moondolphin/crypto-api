package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/app"
)

type GetCurrentPriceHandler struct {
	UC app.GetCurrentPriceUseCase
}

// @Summary Consultar precio actual de criptomoneda
// @Description Consulta el precio actual de una criptomoneda habilitada
// @Tags Crypto
// @Param symbol query string true "Símbolo (BTC, ETH)"
// @Param currency query string true "Moneda (USD, USDT)"
// @Param provider query string true "Proveedor (binance, coingecko)"
// @Success 200 {object} domain.PriceQuote
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/v1/crypto/price [get]
func (h GetCurrentPriceHandler) Handle(c *gin.Context) {
	in := app.GetCurrentPriceInput{
		Symbol:   c.Query("symbol"),
		Currency: c.Query("currency"),
		Provider: c.Query("provider"),
	}

	out, err := h.UC.Execute(c.Request.Context(), in)
	if err != nil {
		switch err {
		case app.ErrBadRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case app.ErrCoinNotEnabled:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case app.ErrProviderNotSupported:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, out)
}
