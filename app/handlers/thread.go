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

type ThreadHandler struct {
	ThreadURL     string
	ThreadUseCase usecases.ThreadUseCase
}

func CreateThreadHandler(router *gin.RouterGroup, threadURL string, threadUseCase usecases.ThreadUseCase) {
	handler := &ThreadHandler{
		ThreadURL:     threadURL,
		ThreadUseCase: threadUseCase,
	}

	threads := router.Group(handler.ThreadURL)
	{
		threads.POST("/:slug_or_id/create", handler.CreatePosts)
		threads.GET("/:slug_or_id/details", handler.GetDetails)
		threads.POST("/:slug_or_id/details", handler.UpdateDetails)
		threads.GET("/:slug_or_id/posts", handler.GetThreadPosts)
		threads.POST("/:slug_or_id/vote", handler.Vote)
	}
}

func (threadHandler *ThreadHandler) CreatePosts(c *gin.Context) {
	slugOrID := c.Param("slug_or_id")

	posts := new(models.Posts)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, posts); err != nil {
		c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
		return
	}

	err := threadHandler.ThreadUseCase.CreatePosts(slugOrID, posts)
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	postsJSON, err := posts.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusCreated, "application/json; charset=utf-8", postsJSON)
}

func (threadHandler *ThreadHandler) GetDetails(c *gin.Context) {
	slugOrID := c.Param("slug_or_id")

	thread, err := threadHandler.ThreadUseCase.Get(slugOrID)
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", threadJSON)
}

func (threadHandler *ThreadHandler) UpdateDetails(c *gin.Context) {
	slugOrID := c.Param("slug_or_id")

	threadUpdate := new(models.ThreadUpdate)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, threadUpdate); err != nil {
		c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
		return
	}

	thread := &models.Thread{
		Title:   threadUpdate.Title,
		Message: threadUpdate.Message,
	}
	err := threadHandler.ThreadUseCase.Update(slugOrID, thread)
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", threadJSON)
}

func (threadHandler *ThreadHandler) GetThreadPosts(c *gin.Context) {
	slugOrID := c.Param("slug_or_id")

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
	sinceStr := c.Query("since")
	since := 0
	if sinceStr != "" {
		var err error
		since, err = strconv.Atoi(sinceStr)
		if err != nil {
			c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
			return
		}
	}
	sort := c.Query("sort")
	if sort == "" {
		sort = "flat"
	}
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

	posts, err := threadHandler.ThreadUseCase.GetPosts(slugOrID, limit, since, sort, desc)
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	postsJSON, err := posts.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", postsJSON)
}

func (threadHandler *ThreadHandler) Vote(c *gin.Context) {
	slugOrID := c.Param("slug_or_id")

	vote := new(models.Vote)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, vote); err != nil {
		c.Data(errors.PrepareErrorResponse(errors.ErrBadRequest))
		return
	}

	thread, err := threadHandler.ThreadUseCase.Vote(slugOrID, vote)
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	threadJSON, err := thread.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", threadJSON)
}
