package ui

import (
	"go-api/cmd/server/middleware"
	"net/http"
	"time"

	cmt "go-api/internal/community"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



func CommunityForm(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("profile")
		if user == nil {
			c.Redirect(http.StatusFound, "/auth/login")
		}
		c.HTML(200, "form-create-community.html", gin.H{ 
			"session_user": user,
		})
}

func CommunitySubmitForm(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	var errors = []gin.H{}

	var community = cmt.Community{}
	if err := c.BindJSON(&community); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}


	community.AuthorID = user.(map[string]interface{})["nickname"].(string)
	community.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	
	if !cmt.RequiredFields(community) {
		c.HTML(200, "response-create-community.html", gin.H{
			"submitted": false,
			"message": "Error creating community",
			"errors": errors,
		})
		return
	}
	

	id, err := cmt.DbCreateCommunity(community)
	if err != nil {
		c.HTML(200, "response-create-community.html", gin.H{
			"submitted": false,
			"message": "Error creating community",
			"errors": errors,
		})
		return
	}

	community.ID = id
	c.HTML(200, "response-create-community.html", gin.H{
		"submitted": true,
		"message": "Community created successfully",
		"errors": errors,
	})
}

func InitializeCommunitiesUI(router *gin.Engine) {
	router.GET("/ui/community/form", middleware.IsAuthenticated, CommunityForm)
	router.POST("/ui/community/submit", middleware.IsAuthenticated, CommunitySubmitForm)
}
