package handlers

import (
	"Technopark_DB_Project/app/models"
	"Technopark_DB_Project/app/usecases"
	"Technopark_DB_Project/pkg/errors"
	"net/http"
	"strconv"

	"github.com/mailru/easyjson"

	"github.com/gin-gonic/gin"
)

type ForumHandler struct {
	ForumURL     string
	ForumUseCase usecases.ForumUseCase
}

func CreateForumHandler(router *gin.RouterGroup, forumURL string, forumUseCase usecases.ForumUseCase) {
	handler := &ForumHandler{
		ForumURL:     forumURL,
		ForumUseCase: forumUseCase,
	}

	forums := router.Group(handler.ForumURL)
	{
		forums.POST("/create", handler.CreateForum)
		forums.GET("/:slug/details", handler.GetDetails)
		forums.POST("/:slug/create", handler.CreateThread)
		forums.GET("/:slug/users", handler.GetForumUsers)
		forums.GET("/:slug/threads", handler.GetForumThreads)
	}
}

func (forumHandler *ForumHandler) CreateForum(c *gin.Context) {
	forum := new(models.Forum)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, forum); err != nil {
		c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
		return
	}

	err := forumHandler.ForumUseCase.CreateForum(forum)
	if err != nil {
		if errors.ResolveErrorToCode(err) == http.StatusConflict {
			forumJSON, err := forum.MarshalJSON()
			if err != nil {
				c.Data(errors.PrepareErrorResponse(err))
				return
			}
			c.Data(errors.ResolveErrorToCode(err), "application/json; charset=utf-8", forumJSON)
		} else {
			c.Data(errors.PrepareErrorResponse(err))
		}
		return
	}

	forumJSON, err := forum.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusCreated, "application/json; charset=utf-8", forumJSON)
}

func (forumHandler *ForumHandler) GetDetails(c *gin.Context) {
	slug := c.Param("slug")

	forum, err := forumHandler.ForumUseCase.Get(slug)
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	forumJSON, err := forum.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", forumJSON)
}

func (forumHandler *ForumHandler) CreateThread(c *gin.Context) {
	slug := c.Param("slug")

	thread := new(models.Thread)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, thread); err != nil {
		c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
		return
	}
	thread.Forum = slug

	err := forumHandler.ForumUseCase.CreateThread(thread)
	if err != nil {
		if errors.ResolveErrorToCode(err) == http.StatusConflict {
			threadJSON, err := thread.MarshalJSON()
			if err != nil {
				c.Data(errors.PrepareErrorResponse(err))
				return
			}
			c.Data(errors.ResolveErrorToCode(err), "application/json; charset=utf-8", threadJSON)
		} else {
			c.Data(errors.PrepareErrorResponse(err))
		}
		return
	}

	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusCreated, "application/json; charset=utf-8", threadJSON)
}

func (forumHandler *ForumHandler) GetForumUsers(c *gin.Context) {
	slug := c.Param("slug")

	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
			return
		}
	}
	since := c.Query("since")
	descStr := c.Query("desc")
	desc := true
	if descStr != "" {
		var err error
		desc, err = strconv.ParseBool(descStr)
		if err != nil {
			c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
			return
		}
	}

	users, err := forumHandler.ForumUseCase.GetUsers(slug, limit, since, desc)
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	usersJSON, err := users.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", usersJSON)
}

func (forumHandler *ForumHandler) GetForumThreads(c *gin.Context) {
	slug := c.Param("slug")

	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
			return
		}
	}
	since := c.Query("since")
	descStr := c.Query("desc")
	desc := true
	if descStr != "" {
		var err error
		desc, err = strconv.ParseBool(descStr)
		if err != nil {
			c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
			return
		}
	}

	threads, err := forumHandler.ForumUseCase.GetThreads(slug, limit, since, desc)
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	threadsJSON, err := threads.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", threadsJSON)
}
