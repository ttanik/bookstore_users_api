package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ttanik/bookstore_oauth-go/oauth"
	"github.com/ttanik/bookstore_oauth-go/oauth/errors"
	"github.com/ttanik/bookstore_users-api/domain/users"
	"github.com/ttanik/bookstore_users-api/services"
	"github.com/ttanik/bookstore_users-api/utils/date_utils"
	"github.com/ttanik/bookstore_utils-go/rest_errors"
)

func Create(c *gin.Context) {
	var user users.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	user.DateCreated = date_utils.GetNowDBFormat()
	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}
func SearchUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implemented")
}
func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}
	userID, userErr := getUserId(c.Param("user_id"))
	if userErr != nil {
		err := errors.NewBadRequestError("invalid user id")
		c.JSON(err.Status, err)
		return
	}
	user, saveErr := services.UsersService.GetUser(userID)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	if oauth.GetCallerId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func Delete(c *gin.Context) {
	userID, userErr := getUserId(c.Param("user_id"))
	if userErr != nil {
		c.JSON(userErr.Status, userErr)
		return
	}
	if err := services.UsersService.DeleteUser(userID); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Update(c *gin.Context) {
	userID, userErr := getUserId(c.Param("user_id"))
	if userErr != nil {
		c.JSON(userErr.Status, userErr)
		return
	}
	var user users.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	user.Id = userID
	isPartial := c.Request.Method == http.MethodPatch

	result, updateErr := services.UsersService.UpdateUser(isPartial, user)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := services.UsersService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users.Marshall(c.GetHeader("X-Public") == "true"))
}

func Login(c *gin.Context) {
	var loginRequest users.LoginRequest
	err := c.ShouldBindJSON(&loginRequest)
	if err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	user, errLogin := services.UsersService.LoginUser(loginRequest)
	if errLogin != nil {
		c.JSON(errLogin.Status, errLogin)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}
func getUserId(idParam string) (int64, *errors.RestErr) {
	userID, userErr := strconv.ParseInt(idParam, 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("invalid user id")
		return 0, err
	}
	return userID, nil
}
