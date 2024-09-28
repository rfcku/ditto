package routes

import (
	"go-api/cmd/server/authenticator"

	ath "go-api/internal/auth"

	"github.com/gin-gonic/gin"
)

func InitializeAuth(router *gin.Engine, a *authenticator.Authenticator) {

	auth := router.Group("/auth")
	{
		auth.GET("/login", ath.Login(a))
		auth.GET("/callback", ath.Callback(a))
		auth.GET("/logout", ath.Logout)
	}
}

