package handler

import (
	"elkonterm/internal/models"
	"elkonterm/internal/repository/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Description Create Category
// @Summary Create Category
// @Tags Category
// @Accept json
// @Produce json
// @Param create body models.CreateCategory true "Create Category"
// @Success 201 {object} models.IdResponse
// @Failure 400,401,404,500 {object} BaseResponse
// @Router /api/v1/admin/categories [post]
// @Security ApiKeyAuth
func (h *Handler) createCategory(c *gin.Context) {
	var body models.CreateCategory

	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.services.Category.CreateCategory(body)
	if err != nil {
		fromError(c, err)
		return
	}

	c.JSON(http.StatusCreated, models.IdResponse{
		ID: id,
	})
}

type getCategoriesResponse struct {
	Categories []models.Category `json:"data"`
	Pagination models.Pagination `json:"meta"`
}

// @Description Get List Category
// @Summary Get List Category
// @Tags Category
// @Accept json
// @Produce json
// @Param page  query int true "page" default(1)
// @Param limit query int true "page limit" default(10)
// @Param search query string false "search key"
// @Param status query bool false "status"
// @Success 200 {object} getCategoriesResponse
// @Failure 400,401,404,500 {object} BaseResponse
// @Router /api/v1/admin/categories [get]
// @Security ApiKeyAuth
func (h *Handler) getListCategory(c *gin.Context) {
	pagination, err := listPagination(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	var options dto.FilterOptions

	filters := make(map[string]any)

	options.Limit = pagination.Limit
	options.Page = pagination.Page

	search := c.Query("search")
	if search != "" {
		filters["search-key"] = search
	}

	status := c.Query("status")
	if status != "" {
		filterStatus, err := strconv.ParseBool(status)
		if err != nil {
			errorResponse(c, http.StatusBadRequest, err)
			return
		}

		filters["status"] = strconv.FormatBool(filterStatus)
	}

	options.Filters = filters

	categories, total, err := h.services.Category.GetList(options)
	if err != nil {
		fromError(c, err)
		return
	}

	updatePagination(&pagination, total)

	c.JSON(http.StatusOK, getCategoriesResponse{
		Categories: categories,
		Pagination: pagination,
	})
}
