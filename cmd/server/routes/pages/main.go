package pages

import (
	"github.com/gin-gonic/gin"
	"go-api/cmd/server/middleware"
	"go-api/internal/community"
	"go-api/internal/post"
	"go-api/internal/user"
)

func InitializePagesRoutes(router *gin.Engine) {

	InitializeHomePage(router)
	InitializeCreatePage(router)

	posts := router.Group("/p")
	{

		posts.GET("/", post.HTMLAllPostsPage)
		posts.GET("/all", post.HTMLAllPostsPage)
		posts.GET("/:id", post.HTMLSinglePostPage)
		posts.GET("/create", middleware.IsAuthenticated, post.HTMLCreatePostPage)

	}

	communities := router.Group("/c")
	{
		communities.GET("/", community.HTMLAllCommunitiesPage)
		communities.GET("/all", community.HTMLAllCommunitiesPage)
		communities.GET("/:name", community.HTMLCommunityPage)
		communities.GET("/create", middleware.IsAuthenticated, community.HTMLCreateCommunityPage)
		communities.GET("/:name/create", middleware.IsAuthenticated, post.HTMLCreatePostPage)

	}

	users := router.Group("/u")
	{
		users.GET("/", user.HTMLAllUsersPage)
		users.GET("/all", user.HTMLAllUsersPage)
		users.GET("/:username", user.HTMLUserPage)
		// users.PUT("/:id", middleware.IsAuthenticated, user.HTMLUpdateUserPage)
	}
}
