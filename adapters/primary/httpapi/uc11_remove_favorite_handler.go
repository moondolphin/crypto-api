package httpapi

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/domain"
)

type RemoveFavoriteHandler struct {
	CoinRepo domain.CoinRepository
	FavRepo  domain.FavoritesRepository
}

// @Summary Quitar moneda de favoritas del usuario
// @Description Elimina la moneda (symbol) de las favoritas del usuario autenticado.
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param symbol path string true "Symbol (BTC, ETH, ...)"
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/v1/users/me/favorites/{symbol} [delete]
func (h RemoveFavoriteHandler) Handle(c *gin.Context) {
	auth, ok := MustAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	symbol := strings.ToUpper(strings.TrimSpace(c.Param("symbol")))
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_symbol"})
		return
	}

	coin, err := h.CoinRepo.GetBySymbol(c.Request.Context(), symbol)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
		return
	}
	if coin == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "coin_not_found"})
		return
	}

	if err := h.FavRepo.RemoveFavoriteCoinFromUser(c.Request.Context(), auth.UserID, coin.ID); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"action": "removed",
		"symbol": symbol,
	})
}
