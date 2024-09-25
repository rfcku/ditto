package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-api/internal/post"
	usr "go-api/internal/user"

	"github.com/gin-contrib/sessions"
)

func InitializeSingleUserPage(router *gin.Engine) {
	router.GET("/u/:username", func (c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		
		username := c.Param("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username is required"})
			return
		}

		profile, err := usr.DbGetUserByUsername(username)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		posts, err := post.DbGetPostsByUser(profile.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		profileView := usr.UserToUserView(profile)
		c.HTML(200, "single-user-page.html", gin.H{
			"session_user": user,
			"title": profileView.Username, 
			"Username": profileView.Username,
			"CreatedAt": profileView.CreatedAt,
			"posts": post.PostsToPostView(posts),
		})
})
}