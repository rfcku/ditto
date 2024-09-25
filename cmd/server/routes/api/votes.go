package api

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/vote"

	"github.com/gin-gonic/gin"
)

func InitializeVotes(router *gin.RouterGroup) {
	router.GET("/", middleware.IsAuthenticated, vote.GetVotes)
	router.POST("/:targetID", middleware.IsAuthenticated,vote.CreateVote)
}