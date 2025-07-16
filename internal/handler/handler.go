package handler

import (
	"github.com/Serhio/backend-vk-test/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// logger and recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.register)
			auth.POST("/login", h.login)
		}

		ads := api.Group("/ads")
		{
			// this route can handle optional authentication
			ads.GET("", h.optionalAuthMiddleware, h.getAds)
			ads.POST("", h.authMiddleware, h.createAd)
		}
	}

	return router
}
