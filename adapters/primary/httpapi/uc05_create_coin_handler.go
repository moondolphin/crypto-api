package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/app"
)

type CreateCoinHandler struct {
	UC app.CreateCoinUseCase
}

// @Summary Alta de moneda de interés
// @Description Crea o actualiza una moneda (upsert). Requiere token.
// @Tags Coins
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body app.CreateCoinInput true "Coin input"
// @Success 200 {object} app.CreateCoinOutput
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/v1/coins [post]
func (h CreateCoinHandler) Handle(c *gin.Context) {
	var in app.CreateCoinInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}
	out, err := h.UC.Execute(c.Request.Context(), in)
	if err != nil {
		switch err {
		case app.ErrInvalidCoinInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case app.ErrCoinNotResolvable:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
		}
		return
	}


	c.JSON(http.StatusOK, out)
}
