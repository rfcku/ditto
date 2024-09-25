package router

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go-api/cmd/server/authenticator"
	"go-api/docs"

	"go-api/cmd/server/routes/api"
	pages "go-api/cmd/server/routes/pages"
	"go-api/cmd/server/routes/ui"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go-api/cmd/server/routes"
)


func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.Next()
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, rate ratelimit.Info) {
	c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
}

func New(auth *authenticator.Authenticator) *gin.Engine {

	router := gin.Default()
	
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{Rate: time.Second, Limit: 15})
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{ErrorHandler: errorHandler, KeyFunc: keyFunc})
	
	gob.Register(map[string]interface{}{})
	
	authStore := cookie.NewStore([]byte("secret"))	

	router.Use(mw)
	router.Use(sessions.Sessions("auth-session", authStore))
	// get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
   	v1 := router.Group("/api/v1")
   	{
      posts := v1.Group("/posts")
      {
		api.InitializePosts(posts)
      }
	  comments := v1.Group("/comments")
	  {
		api.InitializeComments(comments)
	  }
	
		communities := v1.Group("/communities")
		{
			api.InitializeCommunities(communities)
		}
	  votes := v1.Group("/votes")
	  {
		api.InitializeVotes(votes)
	  }
	  awards := v1.Group("/awards")
	  {
		api.InitializeAwards(awards)
	  }
   	}

	server_dir := filepath.Join(dir,"cmd", "server", "html")
	
	assets_dir := filepath.Join(server_dir, "assets")
	templates_dir := filepath.Join(server_dir,"templates", "/**/*")

	// serve the swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
  
	// serve the static files
	router.Static("/assets", assets_dir)
	router.LoadHTMLGlob(templates_dir)

	// initialize the routes
	routes.InitializeAuth(router, auth)
	
	// UI routes - Rendered HTML
	ui.InitializeCommentsUI(router)
	ui.InitializeVotesUI(router)
	ui.InitializePostsUI(router)
	ui.InitializeCommunitiesUI(router)
	// Pages routes
	pages.InitializeHomePage(router)
	pages.InitializeCreatePage(router)
	pages.InitializeSinglePostPage(router)
	pages.InitializeSingleUserPage(router)




	return router
}
