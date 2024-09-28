package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"

	pr "go-api/internal/post"
	"go-api/pkg/utils"

	"github.com/gin-contrib/sessions"
)

func InitializeHomePage(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("profile")

		page, limit, sortBy := utils.DefaultPaginationQueryParams(c)
		posts, total, err := pr.DbGetAllPosts(page, limit, sortBy, user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		paginatedPosts := pr.ToPostView(posts)
		pagination := utils.BuildPagination(page, limit, sortBy, total)

		c.HTML(200, "home-page.html", gin.H{
			"session_user": user,
			"posts":        paginatedPosts,
			"pagination":   pagination,
		})
	})
}
