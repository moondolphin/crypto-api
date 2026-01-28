package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/app"
)

type RefreshHandler struct {
	UC app.RefreshQuotesUseCase
}

// @Summary Refresh quotes (Internal job)
// @Description Consulta providers externos para todas las coins habilitadas y persiste en BD
// @Tags Job
// @Produce json
// @Success 200 {object} app.RefreshQuotesOutput
// @Failure 401 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /api/v1/job/refresh [post]
func (h RefreshHandler) Handle(c *gin.Context) {
	out, err := h.UC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
		return
	}
	c.JSON(http.StatusOK, out)
}
