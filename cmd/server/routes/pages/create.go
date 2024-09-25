package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-api/cmd/server/middleware"

	"github.com/gin-contrib/sessions"

	usr "go-api/internal/user"
	wlt "go-api/internal/wallet"
)

func InitializeCreatePage(router *gin.Engine) {
	router.GET("/create", middleware.IsAuthenticated, func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("profile")
		if user == nil {
			c.Redirect(http.StatusFound, "/auth/login")
		}
		nickname := usr.UserNickName(user)
		var balance int
		b, err := wlt.DbGetUserWalletBalanceByNickName(nickname)
		if err != nil {
			balance = 0
		} else {
			balance = b
		}
		c.HTML(200, "create-post-page.html", gin.H{
			"title": "Submit Post", 
			"session_user": user,
			"balance": balance,
		})
	})
	
	router.GET("/create/community", middleware.IsAuthenticated, func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("profile")
		if user == nil {
			c.Redirect(http.StatusFound, "/auth/login")
		}
		nickname := usr.UserNickName(user)
		
		var balance int
		b, err := wlt.DbGetUserWalletBalanceByNickName(nickname)
		if err != nil {
			balance = 0
		} else {
			balance = b
		}
		c.HTML(200, "create-community-page.html", gin.H{
			"title": "Submit Post", 
			"session_user": user,
			"balance": balance,
		})
	})
}