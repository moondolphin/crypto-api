package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/app"
	"github.com/moondolphin/crypto-api/domain"
)

var _ = domain.PriceQuote{} // para que compile el import (swagger usa domain.PriceQuote)

// GET guarda mismo nombre de handler, pero ahora lee de BD
type GetCurrentPriceHandler struct {
	UC app.GetLastPriceUseCase
}

// @Summary Consultar precio actual de criptomoneda (desde BD)
// @Description Devuelve la última cotización guardada en la base (no consulta al exchange)
// @Tags Crypto
// @Param symbol query string true "Símbolo (BTC, ETH)"
// @Param currency query string false "Moneda (USD, USDT) - opcional"
// @Param provider query string false "Proveedor (binance, coingecko) - opcional"
// @Success 200 {object} domain.PriceQuote
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/v1/crypto/price [get]
func (h GetCurrentPriceHandler) Handle(c *gin.Context) {
	in := app.GetLastPriceInput{
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
		case app.ErrQuoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, out)
}
