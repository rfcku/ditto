package api

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/community"

	"github.com/gin-gonic/gin"
)

func InitializeCommunities(router *gin.RouterGroup) {
	router.GET("/", middleware.IsAuthenticated, community.GetCommunities)
	router.GET("/:id", middleware.IsAuthenticated, community.GetCommunityByID)
	router.POST("/", middleware.IsAuthenticated, community.CreateCommunity)
	router.DELETE("/:id", middleware.IsAuthenticated, community.DeleteCommunity)
}
