package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/app"
)

type RefreshHandler struct {
	UC app.ManualRefreshWithCooldownUseCase
}

// @Summary Refresh quotes (Manual, with cooldown)
// @Description Refresh manual: consulta providers y persiste en BD. Requiere JWT. Enforce cooldown (ej 20 min).
// @Tags Job
// @Produce json
// @Security BearerAuth
// @Success 200 {object} app.ManualRefreshWithCooldownOutput
// @Failure 401 {object} map[string]string
// @Failure 429 {object} map[string]any
// @Failure 503 {object} map[string]string
// @Router /api/v1/job/refresh [post]
func (h RefreshHandler) Handle(c *gin.Context) {
	out, err := h.UC.Execute(c.Request.Context())
	if err != nil {
		switch err {
		case app.ErrCooldownActive:
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":               "cooldown_active",
				"retry_after_seconds": out.RetryAfterSeconds,
			})
			return
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
			return
		}
	}
	c.JSON(http.StatusOK, out)
}
