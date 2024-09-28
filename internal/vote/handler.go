package vote

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @BasePath /api/v1

// GetVotes godoc
// @Summary Get all votes
// @Schemes
// @Description Get all votes
// @Tags votes
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit number"
// @Success 200 {object} Vote
// @Router /votes [get]
func GetVotes(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")

	votes, err := DbGetAllVotes(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, votes)
}

// GetVoteByID godoc
// @Summary Get a vote by ID
// @Schemes
// @Description Get a vote by ID
// @Tags votes
// @Accept json
// @Produce json
// @Param id path string true "Vote ID"
// @Success 200 {object} Vote
// @Router /votes/:id [get]
func GetVoteByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	vote, err := DbGetVoteID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Vote not found"})
		return
	}
	c.JSON(http.StatusOK, vote)
}

// CreateVote godoc
// @Summary Create a vote
// @Schemes
// @Description Create a vote
// @Tags votes
// @Accept json
// @Produce json
// @Param targetID path string true "Target ID"
// @Success 201 {object} Vote
// @Router /votes/{targetID} [post]
func CreateVote(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	postId := c.Param("targetID")
	id, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid vote ID"})
		return
	}

	votes, voted, err := DbSubmitVote("p", id, user.(map[string]interface{})["nickname"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"ID": postId, "Votes": votes, "Voted": voted})
}

// UpdateVote godoc
// @Summary Update a vote
// @Schemes
// @Description Update a vote
// @Tags votes
// @Accept json
// @Produce json
// @Param id path string true "Vote ID"
// @Success 200 {object} Vote
// @Router /votes/:id [put]
func UpdateVote(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	var updatedVote Vote
	if err := c.BindJSON(&updatedVote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = DbUpdateVote(id, updatedVote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	updatedVote.ID = id
	c.JSON(http.StatusOK, updatedVote)
}

// DeleteVote godoc
// @Summary Delete a vote
// @Schemes
// @Description Delete a vote
// @Tags votes
// @Accept json
// @Produce json
// @Param id path string true "Vote ID"
// @Success 200 {object} Vote
// @Router /votes/:id [delete]
func DeleteVote(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeleteVote(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vote deleted"})
}

// BasePath /ui/votes

// HTMLSubmitVote godoc
// @Summary Submit a vote
// @Description Submit a vote
// @Tags votes
// @Accept json
// @Produce html
// @Param targetID path string true "Target ID"
// @Param t query string true "Target type"
// @Success 201 {string} string "HTML Form"
// @Router /html/:targetID [post]
func HTMLSubmitVote(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	targetID := c.Param("targetID")
	id, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid vote ID"})
		return
	}

	target := c.Query("t")
	if target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "target type is required"})
		return
	}

	if target != "p" && target != "c" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid target type"})
		return
	}

	nickname := user.(map[string]interface{})["nickname"].(string)
	votes, voted, err := DbSubmitVote(target, id, nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var template string
	if target == "p" {
		template = "component-vote-post.html"
	} else {
		template = "component-vote-comment.html"
	}

	c.HTML(201, template, gin.H{
		"ID":         targetID,
		"VotesTotal": votes,
		"Voted":      voted,
	})
}
