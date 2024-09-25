package community

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


// @BasePath /api/v1

// GetCommunities godoc
// @Summary Get a community
// @Schemes 
// @Description Get Communities
// @Tags community
// @Accept json
// @Produce json
// @Success 200 {object} Community
// @Router /community/ [get]
func GetCommunities(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := CommunitiesDefaultQueryParams(c)
	communitys, err := DbGetAllCommunities(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, communitys)
}

// @BasePath /api/v1

// GetCommunityByID godoc
// @Summary Get a community by ID
// @Schemes
// @Description Get a community by ID
// @Tags community
// @Accept json
// @Produce json
// @Success 200 {object} Community
// @Router /community/:id [get]
func GetCommunityByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	community, err := DbGetCommunityID(id)
	if err != nil {
		fmt.Println("Error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"message": "Community not found"})
		return
	}
	c.JSON(http.StatusOK, community)
}


// CreateCommunity godoc
// @Summary Create a community
// @Schemes
// @Description Create a community
// @Tags community
// @Accept json
// @Produce json
// @Success 201 {object} Community
// @Router /community [post]
func CreateCommunity(c *gin.Context) {

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

	communityTypeID := c.Param("typeID")
	typeID, err := primitive.ObjectIDFromHex(communityTypeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid type ID"})
		return
	}

	_, err = DbCommunityTypeExists(typeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var community = Community{}
	
	community.AuthorID = user.(map[string]interface{})["id"].(string)
	community.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = DbCreateCommunity(community)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	community.ID = id
	c.JSON(http.StatusCreated, community)
}


// DeleteCommunity godoc
// @Summary Delete a community
// @Schemes
// @Description Delete a community
// @Tags community
// @Accept json
// @Produce json
// @Success 200 {object} Community
// @Router /community/:id [delete]
func DeleteCommunity(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeleteCommunity(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Community deleted"})
}
