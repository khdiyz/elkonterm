package handler

import (
	"elkonterm/config"
	"elkonterm/internal/models"
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

func listPagination(c *gin.Context) (models.Pagination, error) {
	page, err := getPageQuery(c)
	if err != nil {
		return models.Pagination{}, err
	}

	limit, err := getLimitQuery(c)
	if err != nil {
		return models.Pagination{}, err
	}

	return models.Pagination{
		Page:  page,
		Limit: limit,
	}, nil
}

func getPageQuery(c *gin.Context) (int, error) {
	pageStr := c.DefaultQuery("page", config.DefaultPage)

	// Convert the page string to an integer
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, fmt.Errorf("invalid page value: %v", err)
	}

	// Validate that the page number is positive
	if page < 1 {
		return 0, fmt.Errorf("page must be greater than or equal to 1")
	}

	return page, nil
}

func getLimitQuery(c *gin.Context) (int, error) {
	limitStr := c.DefaultQuery("limit", config.DefaultLimit)

	// Convert the limit string to an integer
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, fmt.Errorf("invalid limit value: %v", err)
	}

	// Validate that the limit is positive
	if limit < 1 {
		return 0, fmt.Errorf("limit must be greater than or equal to 1")
	}

	return limit, nil
}

func updatePagination(pagination *models.Pagination, totalItems int) {
	if pagination == nil {
		return
	}

	// Calculate the total number of pages
	pageCount := int(math.Ceil(float64(totalItems) / float64(pagination.Limit)))

	// Update the Pagination object
	pagination.PageCount = pageCount
	pagination.TotalCount = totalItems
}
