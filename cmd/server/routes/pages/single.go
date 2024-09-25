package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	cm "go-api/internal/comment"
	pr "go-api/internal/post"

	"github.com/gin-contrib/sessions"
)

func InitializeSinglePostPage(router *gin.Engine) {
	router.GET("/:targetID", func (c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("profile")
		
		targetID := c.Param("targetID")
		id, err := primitive.ObjectIDFromHex(targetID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
			return
		}

		post, err := pr.DbGetPostID(id, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		comments, err := cm.DbGetAllComments(1, 10, "best", user, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		postView := pr.PostToPostView(post)
		c.HTML(200, "single-post-page.html", gin.H{
			"title": "syntax error", 
			"ID": targetID,
			"TargetID": targetID,
			"session_user": user,
			"Title": postView.Title,
			"Content": postView.Content,
			"AuthorID": postView.AuthorID,
			"CreatedAt": postView.CreatedAt,
			"VotesTotal": postView.VotesTotal,
			"Voted": postView.Voted,
			"CommentsTotal": postView.CommentsTotal,
			"Awards": postView.Awards,
			"AwardsTotal": postView.AwardsTotal,
			"Tags": postView.Tags,
			"Comments": comments,
		})
})
}