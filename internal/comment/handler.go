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
