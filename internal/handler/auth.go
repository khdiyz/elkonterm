package handler

import (
	"elkonterm/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Description Login Admin
// @Summary Login Admin
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login Admin"
// @Success 200 {object} models.LoginResponse
// @Failure 400,401,404,500 {object} BaseResponse
// @Router /api/v1/admin/auth/login [post]
func (h *Handler) loginAdmin(c *gin.Context) {
	var body models.LoginRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	accessToken, refreshToken, err := h.services.Authorization.LoginAdmin(body)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	})
}
