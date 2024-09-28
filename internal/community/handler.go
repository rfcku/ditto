package community

import (
	"fmt"
	"net/http"
	"time"

	"go-api/pkg/utils"

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

	page, limit, sortBy := utils.DefaultPaginationQueryParams(c)
	communities, total, err := DbGetAllCommunities(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := utils.PaginatedResponse{
		Data:       communities,
		Pagination: utils.BuildPagination(page, limit, sortBy, total),
	}

	c.JSON(http.StatusOK, response)
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

// HTMLAllCommunities godoc
// @Summary HTML Community HTMLAllCommunities
// @Schemes
// @Description HTML Community HTMLAllCommunities
// @Tags community
// @Accept json
// @Produce html
// @Success 200 {object} Community
// @Router /ui/community [get]
func HTMLAllCommunities(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := utils.DefaultPaginationQueryParams(c)
	communities, total, err := DbGetAllCommunities(page, limit, sortBy, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	communitiesView := ToCommunityView(communities)
	pagination := utils.BuildPagination(page, limit, sortBy, total)

	c.HTML(200, "list-communities.html", gin.H{
		"session_user": user,
		"communities":  communitiesView,
		"pagination":   pagination,
	})

}

func HTMLCommunityByID(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	targetID := c.Param("targetID")
	id, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	community, err := DbGetCommunityID(id)
	if err != nil {
		c.HTML(404, "404.html", gin.H{})
		return
	}

	c.HTML(200, "single-community.html", gin.H{
		"ID":           targetID,
		"session_user": user,
		"name":         community.Name,
	})

}

// HTMLCommunityForm godoc
// @Summary HTML Community HTMLCommunityForm
// @Schemes
// @Description HTML Community HTMLCommunityForm
// @Tags community
// @Accept json
// @Produce html
// @Success 200 {object} Community
// @Router /ui/community/form [get]
func HTMLCommunityForm(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}
	c.HTML(200, "form-create-community.html", gin.H{
		"session_user": user,
	})
}

// HTMLCommunitySubmitForm godoc
// @Summary HTML Community HTMLCommunitySubmitForm
// @Schemes
// @Description HTML Community HTMLCommunitySubmitForm
// @Tags community
// @Accept json
// @Produce html
// @Success 200 {object} Community
// @Router /ui/community/form [post]
func HTMLSubmitCommunity(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.HTML(http.StatusUnauthorized, "response-create-community.html", gin.H{
			"message": "Error creating community",
			"errors":  "Unauthorized",
		})
		return
	}
	var community = Community{}
	if err := c.BindJSON(&community); err != nil {
		c.HTML(http.StatusBadRequest, "response-create-community.html", gin.H{
			"message": "Error creating community",
			"errors":  err.Error(),
		})
		return
	}

	community.AuthorID = user.(map[string]interface{})["nickname"].(string)
	community.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	if err := utils.ValidateObject(community); err != nil {
		c.HTML(http.StatusBadRequest, "response-create-community.html", gin.H{
			"message": "Error creating community",
			"errors":  err.Error(),
		})
		return
	}

	_, err := DbCreateCommunity(community)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "response-create-community.html", gin.H{
			"message": "Error creating community",
			"errors":  err.Error(),
		})
		return
	}

	c.HTML(http.StatusAccepted, "response-create-community.html", gin.H{
		"submitted": true,
		"message":   "Community created successfully",
	})
}

// HTMLSingleCommunity godoc
// @Summary HTML Community HTMLSingleCommunity
// @Schemes
// @Description HTML Community HTMLSingleCommunity
// @Tags community
// @Accept json
// @Produce html
// @Param targetID path string true "Target ID"
// @Success 200 {object} Community
// @Router /ui/community/:targetID [get]
func HTMLSingleCommunity(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	targetID := c.Param("targetID")
	id, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		c.HTML(http.StatusBadRequest, "page-single-community.html", gin.H{
			"ID":           targetID,
			"session_user": user,
			"error":        err.Error(),
		})
		return
	}

	community, err := DbGetCommunityID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "page-single-community.html", gin.H{
			"ID":           targetID,
			"session_user": user,
			"error":        err.Error(),
		})
		return
	}

	c.HTML(http.StatusAccepted, "page-single-community.html", gin.H{
		"ID":           targetID,
		"session_user": user,
		"name":         community.Name,
	})
}

// HTMLAllCommunitiesPage godoc
// @Summary HTML Community HTMLAllCommunitiesPage
// @Schemes
// @Description HTML Community HTMLAllCommunitiesPage
// @Tags community
// @Accept json
// @Produce html
// @Success 200 {object} Community
// @Router /ui/community [get]
func HTMLAllCommunitiesPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := utils.DefaultPaginationQueryParams(c)
	communities, total, err := DbGetAllCommunities(page, limit, sortBy, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	communitiesView := ToCommunityView(communities)
	pagination := utils.BuildPagination(page, limit, sortBy, total)
	c.HTML(200, "page-communities.html", gin.H{
		"session_user": user,
		"communities":  communitiesView,
		"pagination":   pagination,
	})
}

// HTMLCommunityPage godoc
// @Summary HTML Community HTMLCommunityPage
// @Schemes
// @Description HTML Community HTMLCommunityPage
// @Tags community
// @Accept json
// @Produce html
// @Success 200 {object} Community
// @Router /ui/community/:targetID [get]
func HTMLCommunityPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	name := c.Param("name")
	community, err := DBGetCommunityByName(name)
	if err != nil {
		c.HTML(404, "404.html", gin.H{
			"message": err.Error(),
		})
	}
	c.HTML(200, "page-single-community.html", gin.H{
		"session_user": user,
		"ID":           community.ID,
		"Name":         community.Name,
	})
}

// HTMLCreateCommunityPage godoc
// @Summary HTML Community HTMLCreateCommunityPage
// @Schemes
// @Description HTML Community HTMLCreateCommunityPage
// @Tags community
// @Accept json
// @Produce html
// @Success 200 {object} Community
// @Router /ui/community/create [get]
func HTMLCreateCommunityPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}

	c.HTML(200, "page-create-community.html", gin.H{
		"title":        "Submit Post",
		"session_user": user,
	})
}

func HTMLSearchCommunities(c *gin.Context) {

	query := c.Request.FormValue("search")
	if query == "" {
		c.HTML(http.StatusBadRequest, "search-community-response.html", gin.H{
			"error": "Query parameter is required",
		})
		return
	}

	communities, err := DbGetSearchCommunities(query)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "search-community-response.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	communitiesView := ToCommunityView(communities)
	c.HTML(200, "search-community-response.html", gin.H{
		"communities": communitiesView,
	})
}
