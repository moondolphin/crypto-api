package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moondolphin/crypto-api/app"
)

type RegisterUserHandler struct {
	UC app.RegisterUserUseCase
}

// @Summary Registrar nuevo usuario
// @Description Crea un nuevo usuario (email único) con password hasheada
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body app.RegisterInput true "Register payload"
// @Success 201 {object} app.UserOutput
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (h RegisterUserHandler) Handle(c *gin.Context) {
	var in app.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}

	out, err := h.UC.Execute(c.Request.Context(), in)
	if err != nil {
		switch err {
		case app.ErrInvalidEmail, app.ErrInvalidPassword, app.ErrInvalidName:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case app.ErrEmailAlreadyRegistered:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "internal_error"})
		}
		return
	}

	c.JSON(http.StatusCreated, out)
}
