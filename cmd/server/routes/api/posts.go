package api

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/post"

	"github.com/gin-gonic/gin"
)

func InitializePosts(router *gin.RouterGroup) {
	router.GET("/", post.GetPosts)
	router.POST("/", middleware.IsAuthenticated, post.CreatePost)
	router.GET("/:id", middleware.IsAuthenticated, post.GetPostByID)
	router.PUT("/:id",  middleware.IsAuthenticated, post.UpdatePost)
	router.DELETE("/:id", middleware.IsAuthenticated, post.DeletePost)
	router.GET("/fake",middleware.IsAuthenticated, post.FakePosts)
	router.POST("/upload", middleware.IsAuthenticated, post.UploadFile)
}