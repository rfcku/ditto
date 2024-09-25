package routes

import (
	"go-api/cmd/server/authenticator"

	"go-api/internal/auth"

	"github.com/gin-gonic/gin"
)


func InitializeAuth(router *gin.Engine, a *authenticator.Authenticator) {
	router.GET("/auth/login", auth.Login(a))
	router.GET("/auth/callback", auth.Callback(a))
	router.GET("/callback", auth.Callback(a))
	router.GET("/auth/logout", auth.Logout)
}