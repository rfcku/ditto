package award

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BasePath /api/v1

// GetAwards godoc
// @Summary Get all awards
// @Description Get all awards
// @Tags awards
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit number"
// @Param sortBy query string false "Sort by"
// @Success 200 {object} Award
// @Router /awards [get]
func GetAwards(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := AwardsDefaultQueryParams(c)
	awards, err := DbGetAllAwards(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, awards)
}

// GetAwardByID godoc
// @Summary Get an award by ID
// @Description Get an award by ID
// @Tags awards
// @Accept json
// @Produce json
// @Param id path string true "Award ID"
// @Success 200 {object} Award
// @Router /awards/:id [get]
func GetAwardByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	award, err := DbGetAwardID(id)
	if err != nil {
		fmt.Println("Error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"message": "Award not found"})
		return
	}
	c.JSON(http.StatusOK, award)
}

// CreateAward godoc
// @Summary Create an award
// @Description Create an award
// @Tags awards
// @Accept json
// @Produce json
// @Param postID path string true "Post ID"
// @Param typeID path string true "Type ID"
// @Success 201 {object} Award
// @Router /awards/:postID/:typeID [post]
func CreateAward(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	postID := c.Param("postID")
	id, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
		return
	}

	awardTypeID := c.Param("typeID")
	typeID, err := primitive.ObjectIDFromHex(awardTypeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid type ID"})
		return
	}

	_, err = DbAwardTypeExists(typeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var award = Award{}
	award.TargetID = id
	award.TypeID = typeID
	award.AuthorID = user.(map[string]interface{})["id"].(string)
	award.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = DbCreateAward(award)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	award.ID = id
	c.JSON(http.StatusCreated, award)
}

// DeleteAward godoc
// @Summary Delete an award
// @Description Delete an award
// @Tags awards
// @Accept json
// @Produce json
// @Param id path string true "Award ID"
// @Success 200 {object} Award
// @Router /awards/:id [delete]
func DeleteAward(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeleteAward(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Award deleted"})
}
