package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"

	pr "go-api/internal/post"
	usr "go-api/internal/user"
	wlt "go-api/internal/wallet"

	"github.com/gin-contrib/sessions"
)

func InitializeHomePage(router *gin.Engine) {
	router.GET("/", func (c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := pr.PostsDefaultQueryParams(c)
	posts,total, err := pr.DbGetAllPosts(page, limit, sortBy, user)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	paginatedPosts := pr.PostPaginatedView(posts, total, page, limit, sortBy)
	nickname := usr.UserNickName(user)

	if nickname == "" {
		c.HTML(200, "home-page.html", gin.H{
			"session_user": user,
			"posts": paginatedPosts,
			"balance": "XXXX",
		})
		return
	}

	balance, err := wlt.DbGetUserWalletBalanceByNickName(nickname)
	if err != nil {
		c.HTML(200, "home-page.html", gin.H{
			"session_user": user,
			"balance": "XXXX",
		})
	}
	c.HTML(200, "home-page.html", gin.H{
		"session_user": user,
		"balance": balance,
		"posts": paginatedPosts,
	})
	})
}
