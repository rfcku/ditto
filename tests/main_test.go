package main

import (
	"encoding/json"
	"go-api/cmd/server"
	"net/http"
	"testing"
	"go-api/internal/post"

	"github.com/alecthomas/assert/v2"
	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)

}

func Test_Posts(t *testing.T) {
	g := gofight.New()
	e := gin.Default()
	server.Run()

	var postID string
	var basePath = "/posts"
	
	
	t.Run("GetPosts", func(t *testing.T) {
		g.GET(basePath).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
	})
	
	t.Run("CreatePost", func(t *testing.T) {
		body := post.FakePost()
		g.POST(basePath).SetBody(body).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
			var p post.Post
			json.Unmarshal(r.Body.Bytes(), &p)
		})
	})

	t.Run("GetPostByID", func(t *testing.T) {
		g.GET(basePath+ "/"+ postID ).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
	})

	t.Run("UpdatePost", func(t *testing.T) {
		body := `{
			"name": "John - Updated",
			"gender": "Does"
		}`
		g.PUT(basePath+"/"+ postID).SetBody(body).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
	})

	t.Run("DeletePost", func(t *testing.T) {
		g.DELETE(basePath + "/"+ postID ).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
	})
}
