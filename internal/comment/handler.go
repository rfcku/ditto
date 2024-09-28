package comment

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BasePath /api/v1

// GetComments godoc
// @Summary Get all comments
// @Description Get all comments
// @Tags comments
// @Accept json
// @Produce json
// @Param targetID path string true "Target ID"
// @Param page query int false "Page number"
// @Param limit query int false "Limit number"
// @Param sortBy query string false "Sort by"
// @Success 200 {object} Comment
// @Router /comments/:targetID [get]
func GetComments(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	targetID := c.Param("targetID")
	id, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid target ID"})
		return
	}

	page, limit, sortBy := CommentsDefaultQueryParams(c)
	comments, err := DbGetAllComments(page, limit, sortBy, user, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)

}

// CreateComment godoc
// @Summary Create a comment
// @Description Create a comment
// @Tags comments
// @Accept json
// @Produce json
// @Param targetID path string true "Target ID"
// @Param comment body Comment true "Comment object"
// @Success 201 {object} Comment
// @Router /comments/:targetID [post]
func CreateComment(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var comment = Comment{}
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	targetID := c.Param("targetID")
	target, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid target ID"})
		return
	}

	comment.TargetID = target
	comment.AuthorID = user.(map[string]interface{})["nickname"].(string)
	comment.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	if !RequiredFields(comment) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields"})
		return
	}

	id, err := DbCreateComment(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	comment.ID = id
	c.JSON(http.StatusCreated, comment)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete a comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {string} string "Comment deleted"
// @Router /comments/:id [delete]
func DeleteComment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeleteComment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}

// HTMLCommentForm godoc
// @Summary HTML Comment HTMLCommentForm
// @Description HTML Comment HTMLCommentForm
// @Tags comments
// @Accept json
// @Produce html
// @Param targetID path string true "Target ID"
// @Success 200 {string} string "HTML Form"
// @Router /comments/form/:targetID [get]
func HTMLCommentForm(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}
	targetID := c.Param("targetID")
	_, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}
	c.HTML(200, "form-create-comment.html", gin.H{
		"session_user": user,
		"TargetID":     targetID,
	})
}

// HTMLCommentSubmitForm godoc
// @Summary HTML Comment HTMLCommentSubmitForm
// @Description HTML Comment HTMLCommentSubmitForm
// @Tags comments
// @Accept json
// @Produce html
// @Param targetID path string true "Target ID"
// @Param comment body Comment true "Comment object"
// @Success 200 {string} string "HTML Form"
// @Router /comments/form/:targetID [post]
func HTMLSubmitcomment(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	var errors = []gin.H{}

	var comment = Comment{}
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	targetID := c.Param("targetID")
	target, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.HTML(200, "response-create-comment.html", gin.H{
			"submitted": false,
			"message":   "Error creating comment",
			"errors":    errors,
		})
	}

	comment.TargetID = target
	comment.AuthorID = user.(map[string]interface{})["nickname"].(string)
	comment.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	if !RequiredFields(comment) {
		c.HTML(200, "response-create-comment.html", gin.H{
			"submitted": false,
			"message":   "Error creating comment",
			"errors":    errors,
		})
		return
	}

	id, err := DbCreateComment(comment)
	if err != nil {
		c.HTML(200, "response-create-comment.html", gin.H{
			"submitted": false,
			"message":   "Error creating comment",
			"errors":    errors,
		})
		return
	}

	comment.ID = id
	c.HTML(200, "response-create-comment.html", gin.H{
		"submitted": true,
		"message":   "Comment created successfully",
		"errors":    errors,
	})
}

// HTMLGetCommentByID godoc
// @Summary HTML Get Comment By HTMLGetCommentByID
// @Description HTML Get Comment By HTMLGetCommentByID
// @Tags comments
// @Accept json
// @Produce html
// @Param commentID path string true "Comment ID"
// @Success 200 {string} string "HTML Form"
// @Router /comments/:commentID [get]
func HTMLGetCommentByID(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}

	commentID := c.Param("commentID")
	id, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	comment, err := DbGetCommentID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	comments, err := DbGetAllComments(1, 10, "best", user, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.HTML(200, "single-comment.html", gin.H{
		"ID":         comment.ID.Hex(),
		"TargetID":   comment.TargetID.Hex(),
		"AuthorID":   comment.AuthorID,
		"Content":    comment.Content,
		"Replies":    comments,
		"CreatedAt":  comment.CreatedAt.Time().Format("2006-01-02 15:04:05"),
		"Voted":      comment.Voted,
		"VotesTotal": comment.VotesTotal,
	})
}
