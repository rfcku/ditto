package routes

import (
	"go-api/cmd/server/middleware"
	"go-api/internal/award"
	"go-api/internal/comment"
	"go-api/internal/community"
	"go-api/internal/post"
	"go-api/internal/user"
	"go-api/internal/vote"

	"github.com/gin-gonic/gin"
)

func InitializeApiRoutes(router *gin.Engine) {

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			posts := v1.Group("/posts")
			{
				posts.GET("/", post.APIGetPosts)
				posts.POST("/", middleware.IsAuthenticated, post.APICreatePost)
				posts.GET("/:id", middleware.IsAuthenticated, post.APIGetPostByID)
				posts.PUT("/:id", middleware.IsAuthenticated, post.APIUpdatePost)
				posts.DELETE("/:id", middleware.IsAuthenticated, post.APIDeletePost)
				posts.GET("/fake", middleware.IsAuthenticated, post.FakePosts)

			}

			comments := v1.Group("/comments")
			{
				comments.GET("/:id", comment.GetComments)
				comments.POST("/:id", middleware.IsAuthenticated, comment.CreateComment)
				comments.DELETE("/:id", middleware.IsAuthenticated, comment.DeleteComment)
			}

			awards := v1.Group("/awards")
			{
				awards.GET("/:id", middleware.IsAuthenticated, award.GetAwardByID)
				awards.POST("/:id/:postID", middleware.IsAuthenticated, award.CreateAward)
				awards.DELETE("/:id", middleware.IsAuthenticated, award.DeleteAward)
			}

			users := v1.Group("/users")
			{
				users.GET("/", middleware.IsAuthenticated, user.GetUsers)
				users.GET("/:id", middleware.IsAuthenticated, user.GetUserByID)
				users.PUT("/:id", middleware.IsAuthenticated, user.UpdateUser)
				users.DELETE("/:id", middleware.IsAuthenticated, user.DeleteUser)
			}

			communities := v1.Group("/communities")
			{
				communities.GET("/", community.GetCommunities)
				communities.POST("/", middleware.IsAuthenticated, community.CreateCommunity)
				communities.GET("/:id", middleware.IsAuthenticated, community.GetCommunityByID)
				communities.DELETE("/:id", middleware.IsAuthenticated, community.DeleteCommunity)
			}
			votes := v1.Group("/votes")
			{
				votes.POST("/", middleware.IsAuthenticated, vote.CreateVote)
				votes.DELETE("/:id", middleware.IsAuthenticated, vote.DeleteVote)
			}

		}
	}

}
