package routes

import (
	"github.com/gin-gonic/gin"
	"go-api/cmd/server/middleware"
	// award "go-api/internal/award"
	cmt "go-api/internal/comment"
	cmm "go-api/internal/community"
	pr "go-api/internal/post"
	// usrs "go-api/internal/user"
	vt "go-api/internal/vote"
)

func InitializeUIRoutes(router *gin.Engine) {

	ui := router.Group("/ui")
	{

		posts := ui.Group("/posts")
		{
			posts.GET("/all", pr.HTMLAllPosts)
			posts.GET("/form", middleware.IsAuthenticated, pr.HTMLPostForm)
			posts.POST("/form/submit", middleware.IsAuthenticated, pr.HTMLSubmitPost)
		}
		communities := ui.Group("/communities")
		{
			communities.GET("/all", cmm.HTMLAllCommunities)
			communities.GET("/form", middleware.IsAuthenticated, cmm.HTMLCommunityForm)
			communities.POST("/form/submit", middleware.IsAuthenticated, cmm.HTMLSubmitCommunity)
			communities.POST("/search", cmm.HTMLSearchCommunities)
		}
		// awards := ui.Group("/awards")
		// {
		// awards.GET("/all", award.GetHTMLAllAwards)
		// awards.GET("/form", middleware.IsAuthenticated, award.GetHTMLSubmitAwardForm)
		// 	awards.POST("/form/submit", middleware.IsAuthenticated, award.HTMLSubmitAward)
		// }
		comments := ui.Group("/comments")
		{
			comments.GET("/form/:targetID", cmt.HTMLCommentForm)
			comments.POST("/:targetID", middleware.IsAuthenticated, cmt.HTMLSubmitcomment)
			// comments.DELETE("/:id", middleware.IsAuthenticated, cmt.HTMLDeleteComment)
			// comment.PUT("/:id", middleware.IsAuthenticated, cmt.HTMLUpdateComment)
		}
		// users := ui.Group("/users")
		// {
		// users.GET("/all", usrs.GetHTMLAllUsers)
		// users.GET("/form", middleware.IsAuthenticated, usrs.GetHTMLSubmitUserForm)
		// users.POST("/form/submit", middleware.IsAuthenticated, usrs.HTMLSubmitUser)
		// users.PUT("/form/submit", middleware.IsAuthenticated, usrs.HTMLUpdateUser)
		// }
		votes := ui.Group("/votes")
		{
			// votes.GET("/all", vt.GetHTMLAllVotes)
			// votes.GET("/form", middleware.IsAuthenticated, vt.GetHTMLSubmitVoteForm)
			votes.POST("/form/submit", middleware.IsAuthenticated, vt.HTMLSubmitVote)
			votes.POST("/submit/:targetID", middleware.IsAuthenticated, vt.HTMLSubmitVote)
		}
	}
}
