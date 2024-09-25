package post

import (
	"encoding/json"
	"fmt"
	usr "go-api/internal/user"
	wlt "go-api/internal/wallet"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FakePost() string {
	post := fakePost()
	return post
}

func FakePosts(c *gin.Context) {

	num := c.Query("num")
	numInt, err := strconv.Atoi(num)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid number"})
		return
	}
	posts := []Post{}
	for i := 0; i < numInt; i++ {
		post := fakePost()
		var p Post 
		_, err := json.Marshal(post)
		if err != nil {
			println(err.Error())
			continue
		}
		err = json.Unmarshal([]byte(post), &p)
		if err != nil {
			println(err.Error())
			continue
		}
		DbCreatePost(p)
	}
	c.JSON(http.StatusCreated, posts)
}

// @BasePath /api/v1

// GetPosts godoc
// @Summary Get all posts
// @Schemes 
// @Description Get all posts
// @Tags posts 
// @Accept json
// @Produce json
// @Success 200 {object} Post
// @Router /posts [get]
func GetPosts(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	
	page, limit, sortBy := PostsDefaultQueryParams(c)
	posts, total, err := DbGetAllPosts(page, limit, sortBy, user)
	
	paginatedPosts := PostPaginatedView(posts, total, page, limit, sortBy)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, paginatedPosts)
}
// GetPostByID godoc
// @Summary Get a post by ID
// @Schemes 
// @Description Get a post by ID
// @Tags posts 
// @Accept json
// @Produce json
// @Success 200 {object} Post
// @Router /posts/:id [get]
func GetPostByID(c *gin.Context) {
	
	session := sessions.Default(c)
	user := session.Get("profile")

	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
		}
	fmt.Println("Getting Post by ID", id, user)

	post, err := DbGetPostID(id, user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
		return
	}
	postView := PostToPostView(post)
	c.JSON(http.StatusOK, postView)
}

// CreatePost godoc
// @Summary Create a new post
// @Schemes 
// @Description Create a new post with the given title, content, link and tags
// @Tags posts 
// @Accept json
// @Produce json
// @Success 200 {object} Post
// @Router /posts [post]
func CreatePost(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var post = PostIncoming{}
	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var newPost = Post{
		Title: post.Title,
		Content: post.Content,
		Link: post.Link,
		Tags: strings.Split(post.Tags, ","),
	}

	newPost.AuthorID = user.(map[string]interface{})["nickname"].(string)
	newPost.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	if !RequiredFields(newPost) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields"})
		return
	}

	id, err := DbCreatePost(newPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newPost.ID = id
	c.JSON(http.StatusCreated, newPost)
}

// UpdatePost godoc
// @Summary Update a post
// @Schemes 
// @Description  Update a post with the given title, content, link and tags
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {object} Post
// @Router /posts/:id [put]
func UpdatePost(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	var updatedPost Post
	if err := c.BindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = DbUpdatePost(id, updatedPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	updatedPost.ID = id
	c.JSON(http.StatusOK, updatedPost)
}

// Delete godoc
// @Summary Delete a post
// @Schemes 
// @Description Delete a post with the given ID
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {object} Post
// @Router /posts/:id [delete]
func DeletePost(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeletePost(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// Random godoc
// @Summary Get a random post
// @Schemes 
// @Description Delete a post with the given ID
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {object} Post
// @Router /posts/random [get]
func GetRandomPost(c *gin.Context) {
	post, err := DbGetRandomPost()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

func GetHTMLCreateForm(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}
	nickname := usr.UserNickName(user)
	balance, err := wlt.DbGetUserWalletBalanceByNickName(nickname)
	if err != nil {
		c.HTML(200, "form-create-post.html", gin.H{
			"title": "Submit Post", 
			"session_user": user,
		})
	}
	c.HTML(200, "form-create-post.html", gin.H{
		"title": "Submit Post", 
		"session_user": user,
		"balance": balance,
	})
}

func GetHTMLAllPosts(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := PostsDefaultQueryParams(c)
	posts, total,  err := DbGetAllPosts(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	paginatedPosts := PostPaginatedView(posts, total, page, limit, sortBy)

	c.HTML(200, "list-post.html", gin.H{
		"session_user": user,
		"posts": paginatedPosts,
	})
}

func GetHTMLSubmitPostForm(c *gin.Context) {
		
	session := sessions.Default(c)
	user := session.Get("profile")

	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}

	var errors = []gin.H{}
	var newPost Post
	if err := c.BindJSON(&newPost); err != nil {
		errors = append(errors, gin.H{"message": err.Error()})
	}
	newPost.AuthorID = user.(map[string]interface{})["nickname"].(string)
	if !RequiredFields(newPost) {
		errors = append(errors, gin.H{"message": "Missing required fields"})
	}

	if len(errors) > 0 {
		c.HTML(200, "response-create-post.html", gin.H{
			"submitted": false,
			"message": "Error creating post",
			"errors": errors,
		})
		return
	}

	id, err := DbCreatePost(newPost)
	if err != nil {
		errors = append(errors, gin.H{"message": err.Error()})
	}
	newPost.ID = id
	c.HTML(200, "response-create-post.html", gin.H{
		"submitted": true,
		"message": "Post created successfully",
		"errors": errors,
	})
}

// UploadFile godoc
// @Summary Upload a file
// @Schemes 
// @Description Upload a file
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {object} Post
// @Router /posts [post]
func UploadFile(c *gin.Context) {
		
		session := sessions.Default(c)
		user := session.Get("profile")
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}
		files := form.File["files"]

		for _, file := range files {
			filename := filepath.Base(file.Filename)
			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
				return
			}
		}

		c.String(http.StatusOK, "Uploaded successfully %d files with fields name=%s and email=%s.", len(files))
}
