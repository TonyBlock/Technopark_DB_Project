package handlers

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/usecases"
	"Technopark_DB_Project/pkg/errors"
	"net/http"

	"github.com/mailru/easyjson"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserURL     string
	UserUseCase usecases.UserUseCase
}

func CreateUserHandler(router *gin.RouterGroup, userURL string, userUseCase usecases.UserUseCase) {
	handler := &UserHandler{
		UserURL:     userURL,
		UserUseCase: userUseCase,
	}

	users := router.Group(handler.UserURL)
	{
		users.POST("/:nickname/create", handler.CreateUser)
		users.GET("/:nickname/profile", handler.GetUser)
		users.POST("/:nickname/profile", handler.UpdateUser)
	}
}

func (userHandler *UserHandler) CreateUser(c *gin.Context) {
	nickname := c.Param("nickname")

	userUpdate := new(models.UserUpdate)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, userUpdate); err != nil {
		c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
		return
	}

	user := &models.User{
		Nickname: nickname,
		Fullname: userUpdate.Fullname,
		About:    userUpdate.About,
		Email:    userUpdate.Email,
	}

	users, err := userHandler.UserUseCase.Create(user)
	if err != nil {
		if errors.ResolveErrorToCode(err) == http.StatusConflict {
			usersJSON, err := users.MarshalJSON()
			if err != nil {
				c.Data(errors.PrepareErrorResponse(err))
				return
			}
			c.Data(errors.ResolveErrorToCode(err), "application/json; charset=utf-8", usersJSON)
		} else {
			c.Data(errors.PrepareErrorResponse(err))
		}
		return
	}

	userJSON, err := user.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusCreated, "application/json; charset=utf-8", userJSON)
}

func (userHandler *UserHandler) GetUser(c *gin.Context) {
	nickname := c.Param("nickname")

	user, err := userHandler.UserUseCase.Get(nickname)
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	userJSON, err := user.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", userJSON)
}

func (userHandler *UserHandler) UpdateUser(c *gin.Context) {
	nickname := c.Param("nickname")

	userUpdate := new(models.UserUpdate)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, userUpdate); err != nil {
		c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
		return
	}

	user := &models.User{
		Nickname: nickname,
		Fullname: userUpdate.Fullname,
		About:    userUpdate.About,
		Email:    userUpdate.Email,
	}

	err := userHandler.UserUseCase.Update(user)
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	userJSON, err := user.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", userJSON)
}
