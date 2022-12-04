package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ttanik/bookstore_users-api/domain/users"
	"github.com/ttanik/bookstore_users-api/services"
	"github.com/ttanik/bookstore_users-api/utils/errors"
)

func CreateUser(c *gin.Context) {
	var user users.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}
func SearchUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implemented")
}
func GetUser(c *gin.Context) {
	userID, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("invalid user id")
		c.JSON(err.Status, err)
		return
	}
	user, saveErr := services.GetUser(userID)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusOK, user)
}
