package api

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/comment"

	"github.com/gin-gonic/gin"
)

func InitializeComments(router *gin.RouterGroup) {
	router.GET("/:targetID", comment.GetComments)
	router.POST("/:targetID", middleware.IsAuthenticated, comment.CreateComment)
	router.DELETE("/:id", middleware.IsAuthenticated, comment.DeleteComment)
}