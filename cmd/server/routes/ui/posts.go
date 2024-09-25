package ui

import (
	"go-api/cmd/server/middleware"
	pr "go-api/internal/post"

	"github.com/gin-gonic/gin"
)

func InitializePostsUI(router *gin.Engine) {
	router.GET("/ui/posts/all", pr.GetHTMLAllPosts)
	router.GET("/ui/posts/form", middleware.IsAuthenticated, pr.GetHTMLSubmitPostForm)
	router.POST("/ui/posts/form/submit", middleware.IsAuthenticated, pr.GetHTMLSubmitPostForm)
}