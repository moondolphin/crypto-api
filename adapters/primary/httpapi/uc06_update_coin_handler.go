package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/app"
	_ "github.com/moondolphin/crypto-api/domain"
)

type UpdateCoinHandler struct {
	UC app.UpdateCoinUseCase
}

// @Summary Actualizar moneda de interés
// @Description Permite habilitar/deshabilitar una moneda o actualizar IDs de proveedores. Requiere token.
// @Tags Coins
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param symbol path string true "Símbolo (BTC, ETH, UNI...)"
// @Param body body app.UpdateCoinInput true "Datos a actualizar (campos opcionales)"
// @Success 200 {object} domain.Coin
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/v1/coins/{symbol} [put]
func (h UpdateCoinHandler) Handle(c *gin.Context) {
	var in app.UpdateCoinInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}
	in.Symbol = c.Param("symbol")

	out, err := h.UC.Execute(c.Request.Context(), in)
	if err != nil {
		switch err {
		case app.ErrInvalidCoinUpdate:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case app.ErrCoinNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
		}
		return
	}

	c.JSON(http.StatusOK, out)
}
