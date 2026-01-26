package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/app"
)

type LoginHandler struct {
	UC app.LoginUseCase
}

// @Summary Login
// @Description Autentica usuario y devuelve token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body app.LoginInput true "Login payload"
// @Success 200 {object} app.LoginOutput
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (h LoginHandler) Handle(c *gin.Context) {
	var in app.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}

	out, err := h.UC.Execute(c.Request.Context(), in)
	if err != nil {
		switch err {
		case app.ErrBadRequest, app.ErrInvalidEmailL, app.ErrInvalidPasswordL:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case app.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
		}
		return
	}

	c.JSON(http.StatusOK, out)
}
