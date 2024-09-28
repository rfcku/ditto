package post

import (
	"encoding/json"
	"fmt"
	cm "go-api/internal/comment"
	cmm "go-api/internal/community"
	fls "go-api/internal/file"
	utils "go-api/pkg/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FakePosts(c *gin.Context) {

	num := c.Query("num")
	numInt, err := strconv.Atoi(num)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid number"})
		return
	}
	posts := []Post{}
	for i := 0; i < numInt; i++ {
		post := utils.FakePost()
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
		_, err = DbCreatePost(p)
		if err != nil {
			println(err.Error())
			continue
		}
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
func APIGetPosts(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := utils.DefaultPaginationQueryParams(c)
	println("Getting Posts", page, limit, sortBy, user)
	posts, total, err := DbGetAllPosts(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := utils.PaginatedResponse{
		Data:       posts,
		Pagination: utils.BuildPagination(page, limit, sortBy, total),
	}
	c.JSON(http.StatusOK, response)
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
func APIGetPostByID(c *gin.Context) {

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
	c.JSON(http.StatusOK, post)
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
func APICreatePost(c *gin.Context) {

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
		Title:   post.Title,
		Content: post.Content,
		Link:    post.Link,
		Tags:    strings.Split(post.Tags, ","),
	}

	newPost.AuthorID = user.(map[string]interface{})["nickname"].(string)
	newPost.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	if err := utils.ValidateObject(newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
func APIUpdatePost(c *gin.Context) {
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
func APIDeletePost(c *gin.Context) {
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
func APIGetRandomPost(c *gin.Context) {
	post, err := DbGetRandomPost()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

// HTMLSinglePost godoc
// @Summary Get the HTML page for a single post
// @Schemes
// @Description Get the HTML element for a single post
// @Tags posts
// @Accept json
// @Produce html
// @Param targetID path string true "Target ID"
// @Success 200 {string} string "HTML Form"
// @Router /posts/:targetID [get]
func HTMLSinglePost(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	targetID := c.Param("targetID")
	id, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	post, err := DbGetPostID(id, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	comments, err := cm.DbGetAllComments(1, 10, "best", user, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	postView := post.View()
	c.HTML(200, "single-post.html", gin.H{
		"title":         "syntax error",
		"ID":            targetID,
		"session_user":  user,
		"Title":         postView.Title,
		"Content":       postView.Content,
		"AuthorID":      postView.AuthorID,
		"CreatedAt":     postView.CreatedAt,
		"VotesTotal":    postView.VotesTotal,
		"Voted":         postView.Voted,
		"CommentsTotal": postView.CommentsTotal,
		"Awards":        postView.Awards,
		"AwardsTotal":   postView.AwardsTotal,
		"Tags":          postView.Tags,
		"Comments":      comments,
	})
}

// GetHTMLCreateForm godoc
// @Summary Get the HTML form to create a post
// @Schemes
// @Description Get the HTML form to create a post
// @Tags posts
// @Accept json
// @Produce html
// @Success 200 {object} Post
// @Router /posts/create [get]
func HTMLPostForm(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	id := c.Param("id")
	if id != "" {

		community, err := cmm.DBGetCommunityByName(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.HTML(200, "form-create-post.html", gin.H{
			"title":        "Submit Post",
			"session_user": user,
			"community":    community,
		})
		return
	}

	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}

	c.HTML(200, "form-create-post.html", gin.H{
		"title":        "Submit Post",
		"session_user": user,
	})
}

// HTMLAllPosts godoc
// @Summary Get the HTML page for all posts
// @Schemes
// @Description Get the HTML list elements for all posts
// @Tags posts
// @Accept json
// @Produce html
// @Success 200 {object} Post
// @Router /posts [get]
func HTMLAllPosts(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := utils.DefaultPaginationQueryParams(c)
	posts, total, err := DbGetAllPosts(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	pagination := utils.BuildPagination(page, limit, sortBy, total)
	postsView := ToPostView(posts)

	c.HTML(200, "list-posts.html", gin.H{
		"session_user": user,
		"posts":        postsView,
		"pagination":   pagination,
	})
}

// GetHTMLSUbmitPost godoc
// @Summary Get the HTML form to submit a post
// @Schemes
// @Description Get the HTML form to submit a post
// @Tags posts
// @Accept json
// @Produce html
// @Success 200 {object} Post
// @Router /posts/submit [post]
func HTMLSubmitPost(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")

	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
		return
	}
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.HTML(200, "response-submit-post.html", gin.H{
			"submitted": false,
			"message":   "Error creating post",
			"error":     err.Error(),
		})
		return
	}

	var savedFiles []primitive.ObjectID

	files := form.File["files"]
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.HTML(200, "response-submit-post.html", gin.H{
				"submitted": false,
				"message":   "Error creating post",
				"error":     err.Error(),
			})
			return
		}

		var fileType = fls.GetFileType(filename)
		var newFile = fls.File{
			AuthorID: user.(map[string]interface{})["nickname"].(string),
			Type:     int8(fileType),
		}

		id, err := fls.DbCreateFile(newFile)
		if err != nil {
			c.HTML(200, "response-submit-post.html", gin.H{
				"submitted": false,
				"message":   "Error creating post",
				"error":     err.Error(),
			})
			return
		}
		savedFiles = append(savedFiles, id)
	}

	// build post object from form data
	var newPost Post
	formData := c.Request.PostForm
	newPost.Title = formData.Get("title")
	newPost.Content = formData.Get("content")
	newPost.Link = formData.Get("link")
	newPost.Tags = strings.Split(formData.Get("tags"), ",")
	newPost.AuthorID = user.(map[string]interface{})["nickname"].(string)
	newPost.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	newPost.Files = savedFiles

	if err := utils.ValidateObject(newPost); err != nil {
		c.HTML(200, "response-submit-post.html", gin.H{
			"submitted": false,
			"message":   "Error creating post",
			"error":     err.Error(),
		})
		return
	}
	id, err := DbCreatePost(newPost)
	if err != nil {
		c.HTML(200, "response-submit-post.html", gin.H{
			"submitted": false,
			"message":   "Error creating post",
			"error":     err.Error(),
		})
		return
	}
	newPost.ID = id
	c.HTML(200, "response-submit-post.html", gin.H{
		"submitted": true,
		"message":   "Post created successfully",
		"error":     "",
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

// HTMLAllPostsPage godoc
// @Summary Get the HTML page for all posts
// @Schemes
// @Description Get the HTML list elements for all posts
// @Tags posts
// @Accept json
// @Produce html
// @Success 200 {object} Post
// @Router /p [get]
func HTMLAllPostsPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := utils.DefaultPaginationQueryParams(c)
	posts, total, err := DbGetAllPosts(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	pagination := utils.BuildPagination(page, limit, sortBy, total)
	postsView := ToPostView(posts)

	c.HTML(200, "home-page.html", gin.H{
		"session_user": user,
		"posts":        postsView,
		"pagination":   pagination,
	})
}

// HTMLSinglePostPage godoc
// @Summary Get the HTML page for a single post
// @Schemes
// @Description Get the HTML element for a single post
// @Tags posts
// @Accept json
// @Produce html
// @Param targetID path string true "Target ID"
// @Success 200 {string} string "HTML Form"
// @Router /p/:id [get]
func HTMLSinglePostPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	targetID := c.Param("id")
	id, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.HTML(404, "404.html", gin.H{})
		return
	}

	post, err := DbGetPostID(id, user)
	if err != nil {
		c.HTML(404, "404.html", gin.H{})
		return
	}

	comments, err := cm.DbGetAllComments(1, 10, "best", user, id)
	if err != nil {
		c.HTML(500, "404.html", gin.H{})
		return
	}

	postView := post.View()
	c.HTML(200, "page-single-post.html", gin.H{
		"ID":            targetID,
		"session_user":  user,
		"Title":         postView.Title,
		"Content":       postView.Content,
		"AuthorID":      postView.AuthorID,
		"CreatedAt":     postView.CreatedAt,
		"VotesTotal":    postView.VotesTotal,
		"Voted":         postView.Voted,
		"CommentsTotal": postView.CommentsTotal,
		"Awards":        postView.Awards,
		"AwardsTotal":   postView.AwardsTotal,
		"Tags":          postView.Tags,
		"Comments":      comments,
		"TargetID":      targetID,
	})
}

// HTMLCreatePostPage godoc
// @Summary Get the HTML page to create a post
// @Schemes
// @Description Get the HTML page to create a post
// @Tags posts
// @Accept json
// @Produce html
// @Success 200 {object} Post
// @Router /p/create [get]
func HTMLCreatePostPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}
	communityID := c.Param("name")

	if communityID != "" && communityID != "0" && communityID != "create" {
		community, err := cmm.DBGetCommunityByName(communityID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.HTML(200, "page-create-post.html", gin.H{
			"title":        "Submit Post",
			"session_user": user,
			"community":    community.Name,
			"communityID":  community.ID,
		})
		return
	}

	c.HTML(200, "page-create-post.html", gin.H{
		"title":        "Submit Post",
		"session_user": user,
	})
}
