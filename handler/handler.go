package handler

import (
	"hospital_track/service"

	"github.com/gin-gonic/gin"
)

type SHandler struct {
	services *service.SService
}

func Handler(services *service.SService) *SHandler {
	return &SHandler{services: services}
}

func (h *SHandler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		api.POST("/tracks", h.tracks)
	}

	return router
}
