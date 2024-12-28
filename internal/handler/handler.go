package handler

import (
	"elkonterm/config"
	"elkonterm/docs"
	"elkonterm/internal/service"
	"elkonterm/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
	logger   *logger.Logger
}

func NewHandler(services *service.Service, logger *logger.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.HandleMethodNotAllowed = true
	router.Use(corsMiddleware())

	// Setup Swagger documentation
	h.setupSwagger(router)

	h.initAdminRoutes(router)

	return router
}

func (h *Handler) initAdminRoutes(router *gin.Engine) {
	router.POST("/api/v1/admin/auth/login", h.loginAdmin)

	admin := router.Group("/api/v1/admin", h.adminIdentity())
	{
		categories := admin.Group("/categories")
		{
			categories.POST("", h.createCategory)
			categories.GET("", h.getListCategory)
		}
	}
}

func (h *Handler) setupSwagger(router *gin.Engine) {
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler), func(ctx *gin.Context) {
		docs.SwaggerInfo.Host = ctx.Request.Host
		if ctx.Request.TLS != nil {
			docs.SwaggerInfo.Schemes = []string{"https"}
		}
	})
}
