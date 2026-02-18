package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/app"
	"github.com/moondolphin/crypto-api/domain"
)

type ListFavoritesHandler struct {
	FavRepo domain.FavoritesRepository
}

// @Summary Listar monedas favoritas del usuario
// @Description Devuelve las coins favoritas del usuario autenticado.
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} app.FavoriteCoinOutput
// @Failure 401 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/v1/users/me/favorites [get]
func (h ListFavoritesHandler) Handle(c *gin.Context) {
	auth, ok := MustAuth(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	coins, err := h.FavRepo.ListFavoriteCoinIDsByUser(c.Request.Context(), auth.UserID)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
		return
	}

	out := make([]app.FavoriteCoinOutput, 0, len(coins))
	for _, coin := range coins {
		out = append(out, app.FavoriteCoinOutput{
			ID:            coin.ID,
			Symbol:        coin.Symbol,
			Enabled:       coin.Enabled,
			CoinGeckoID:   coin.CoinGeckoID,
			BinanceSymbol: coin.BinanceSymbol,
		})
	}

	c.JSON(http.StatusOK, out)
}
