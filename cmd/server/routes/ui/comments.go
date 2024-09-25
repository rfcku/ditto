package ui

import (
	"go-api/cmd/server/middleware"
	"net/http"
	"time"

	cmt "go-api/internal/comment"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



func CommentForm(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("profile")
		if user == nil {
			c.Redirect(http.StatusFound, "/auth/login")
		}
		targetID := c.Param("targetID")
		_, err := primitive.ObjectIDFromHex(targetID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
			return
		}
		c.HTML(200, "form-create-comment.html", gin.H{ 
			"session_user": user,
			"TargetID": targetID,
		})
}

func CommentSubmitForm(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	var errors = []gin.H{}

	var comment = cmt.Comment{}
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	targetID := c.Param("targetID")
	target, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.HTML(200, "response-create-comment.html", gin.H{
			"submitted": false,
			"message": "Error creating comment",
			"errors": errors,
		})
	}

	comment.TargetID = target
	comment.AuthorID = user.(map[string]interface{})["nickname"].(string)
	comment.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	
	if !cmt.RequiredFields(comment) {
		c.HTML(200, "response-create-comment.html", gin.H{
			"submitted": false,
			"message": "Error creating comment",
			"errors": errors,
		})
		return
	}
	

	id, err := cmt.DbCreateComment(comment)
	if err != nil {
		c.HTML(200, "response-create-comment.html", gin.H{
			"submitted": false,
			"message": "Error creating comment",
			"errors": errors,
		})
		return
	}

	comment.ID = id
	c.HTML(200, "response-create-comment.html", gin.H{
		"submitted": true,
		"message": "Comment created successfully",
		"errors": errors,
	})
}

func InitializeCommentsUI(router *gin.Engine) {
	router.GET("/ui/comment/:commentID", middleware.IsAuthenticated, func(c *gin.Context) {
		
		session := sessions.Default(c)
		user := session.Get("profile")
		if user == nil {
			c.Redirect(http.StatusFound, "/auth/login")
		}

		commentID := c.Param("commentID")
		id, err := primitive.ObjectIDFromHex(commentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
			return
		}

		comment, err := cmt.DbGetCommentID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		comments, err := cmt.DbGetAllComments(1, 10, "best", user, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		
		c.HTML(200, "single-comment.html", gin.H{
			"ID": comment.ID.Hex(),
			"TargetID": comment.TargetID.Hex(),
			"AuthorID": comment.AuthorID,
			"Content": comment.Content,
			"Replies": comments,
			"CreatedAt": comment.CreatedAt.Time().Format("2006-01-02 15:04:05"),
			"Voted": comment.Voted,
			"VotesTotal": comment.VotesTotal,
		})
	})
	router.GET("/ui/comment/form/:targetID", middleware.IsAuthenticated, CommentForm)
	router.POST("/ui/comment/submit/:targetID", middleware.IsAuthenticated, CommentSubmitForm)
}