package user

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @BasePath /api/v1

// GetUsers godoc
// @Summary Get all users
// @Schemes
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Router /users/ [get]
func GetUsers(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")

	users, err := DbGetAllUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}


// GetUserByID godoc
// @Summary Get a user by ID
// @Schemes
// @Description Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Router /users/:id [get]
func GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	user, err := DbGetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}


// CreateUser godoc
// @Summary Create a user
// @Schemes
// @Description Create a user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Router /users [post]
func CreateUser(c *gin.Context) {
	
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := DbCreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}


// UpdateUser godoc
// @Summary Update a user
// @Schemes
// @Description Update a user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Router /users/:id [put]
func UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	var updatedUser User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = DbUpdateUser(id, updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	updatedUser.ID = id
	c.JSON(http.StatusOK, updatedUser)
}


// DeleteUser godoc
// @Summary Delete a user
// @Schemes
// @Description Delete a user with the given ID
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Router /users/:id [delete]
func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
