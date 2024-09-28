package pages

import (
	"github.com/gin-gonic/gin"
	"go-api/cmd/server/middleware"

	"go-api/internal/community"
	"go-api/internal/post"
)

func InitializeCreatePage(router *gin.Engine) {

	router.GET("/create", middleware.IsAuthenticated, post.HTMLCreatePostPage)

	router.GET("/create/community", middleware.IsAuthenticated, community.HTMLCreateCommunityPage)

}
