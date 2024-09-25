package ui

import (
	"net/http"

	"go-api/cmd/server/middleware"
	vt "go-api/internal/vote"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-contrib/sessions"
)

func UiSubmitVote(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	targetID := c.Param("targetID")
	id, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid vote ID"})
		return
	}
	
	target := c.Query("t")
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "target type is required"})
		return
	}

	if target != "p" && target != "c" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid target type"})
		return
	}

	nickname := user.(map[string]interface{})["nickname"].(string)
	votes, voted, err := vt.DbSubmitVote(target, id, nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var template string
	if target == "p" {
		template = "component-vote-post.html"
	} else {
		template = "component-vote-comment.html"
	}

	c.HTML(201, template, gin.H{
		"ID":         targetID,
		"VotesTotal": votes,
		"Voted":      voted,
	})
}

func InitializeVotesUI(router *gin.Engine) {
	router.POST("/ui/vote/submit/:targetID", middleware.IsAuthenticated, UiSubmitVote)
}

