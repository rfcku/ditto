package file

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BasePath /api/v1

// GetPostFiles godoc
// @Summary Get all files from post
// @Description Get all files from post
// @Tags files
// @Accept json
// @Produce json
// @Param postID path string true "Post ID"
// @Param page query int false "Page number"
// @Param limit query int false "Limit number"
// @Param sortBy query string false "Sort by"
// @Success 200 {object} File
// @Router /files/:postID [post]
func Create(c *gin.Context) {

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

	fileTypeID := c.Param("typeID")
	typeID, err := primitive.ObjectIDFromHex(fileTypeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid type ID"})
		return
	}

	var file = File{}
	file.TargetID = id
	file.TypeID = typeID
	file.AuthorID = user.(map[string]interface{})["id"].(string)
	file.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = DbCreateFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	file.ID = id
	c.JSON(http.StatusCreated, file)
}

// Read godoc
// @Summary Get a file by ID
// @Description Get a file by ID
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} File
// @Router /files/:id/read/ [get]
func Read(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	file, err := DbGetFileID(id)
	if err != nil {
		fmt.Println("Error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"message": "File not found"})
		return
	}
	c.JSON(http.StatusOK, file)
}

// Update godoc
// @Summary Update a file
// @Description Update a file with the given title, content, link and tags
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} File
// @Router /files/:id/update [put]
func Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	var file = File{}
	if err := c.BindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	file.ID = id
	err = DbUpdateFile(file.ID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, file)
}

// Delete godoc
// @Summary Delete a file
// @Description Delete a file with the given ID
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} File
// @Router /files/:id [delete]
func Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	err = DbDeleteFile(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "File deleted"})
}

func GetFiles(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := FilesDefaultQueryParams(c)
	files, err := DbGetAllFiles(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

// GetFileByID godoc
// @Summary Get a file by ID
// @Description Get a file by ID
// @Tags files
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} File
// @Router /files/:id [get]
func GetFileByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	file, err := DbGetFileID(id)
	if err != nil {
		fmt.Println("Error", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"message": "File not found"})
		return
	}
	c.JSON(http.StatusOK, file)
}
