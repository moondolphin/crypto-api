package httpapi

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/domain"
)

type AddFavoriteHandler struct {
	CoinRepo domain.CoinRepository
	FavRepo  domain.FavoritesRepository
}

// @Summary Agregar moneda a favoritas del usuario
// @Description Agrega la moneda (symbol) a las favoritas del usuario autenticado.
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param symbol path string true "Symbol (BTC, ETH, ...)"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/v1/users/me/favorites/{symbol} [post]
func (h AddFavoriteHandler) Handle(c *gin.Context) {
	auth, ok := MustAuth(c)
	if !ok {
		// si AuthRequired está puesto, esto no debería pasar
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

	if err := h.FavRepo.AddFavoriteCoinToUser(c.Request.Context(), auth.UserID, coin.ID); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"action": "added",
		"symbol": symbol,
	})
}
