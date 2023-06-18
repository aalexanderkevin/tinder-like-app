package handler

import (
	"net/http"
	"tinder-like-app/container"
	"tinder-like-app/controller/web/request"
	"tinder-like-app/controller/web/response"
	"tinder-like-app/helper"
	"tinder-like-app/model"
	"tinder-like-app/usecase"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type User struct {
	appContainer *container.Container
}

func NewUser(appContainer *container.Container) *User {
	return &User{appContainer: appContainer}
}

func (u *User) SignUp(c *gin.Context) {
	logger := helper.GetLogger(c).WithField("method", "Web.User.SignUp")

	var userRequest request.User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		logger.WithError(err).Warning("Bad request")
		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
		return
	}
	if err := userRequest.Validate(); err != nil {
		response.SendErrorResponse(c, response.ErrValidation, "")
		return
	}

	passwordSalt := ksuid.New().String()
	newUser := model.User{
		Name:         userRequest.Name,
		Email:        userRequest.Email,
		Password:     helper.Pointer(helper.Hash(passwordSalt, *userRequest.Password)),
		PasswordSalt: helper.Pointer(passwordSalt),
	}

	userUseCase := usecase.NewUser(u.appContainer)
	createdUser, err := userUseCase.SignUp(c, newUser)
	if err != nil {
		logger.Error(err.Error())
		if helper.IsNotFoundError(err) {
			response.SendErrorResponse(c, response.ErrDuplicate, err.Error())
			return
		}
		response.SendErrorResponse(c, response.ErrServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (u *User) Login(c *gin.Context) {
	logger := helper.GetLogger(c).WithField("method", "Web.User.Login")

	var userRequest request.User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		logger.WithError(err).Warning("Bad request")
		response.SendErrorResponse(c, response.ErrBadRequest, err.Error())
		return
	}
	if err := userRequest.ValidateLogin(); err != nil {
		response.SendErrorResponse(c, response.ErrValidation, "")
		return
	}

	newUser := model.User{
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}

	userUseCase := usecase.NewUser(u.appContainer)
	_, ok := userUseCase.Login(c, newUser)
	if !ok {
		response.SendErrorResponse(c, response.ErrUnauthorized, "")
		return
	}

	c.JSON(http.StatusOK, "Success Login")
}
