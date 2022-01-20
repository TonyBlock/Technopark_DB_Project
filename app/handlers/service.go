package handlers

import (
	"Technopark_DB_Project/app/usecases"
	"Technopark_DB_Project/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	ServiceURL     string
	ServiceUseCase usecases.ServiceUseCase
}

func CreateServiceHandler(router *gin.RouterGroup, serviceURL string, serviceUseCase usecases.ServiceUseCase) {
	handler := &ServiceHandler{
		ServiceURL:     serviceURL,
		ServiceUseCase: serviceUseCase,
	}

	service := router.Group(handler.ServiceURL)
	{
		service.POST("/clear", handler.Clear)
		service.GET("/status", handler.GetStatus)
	}
}

func (serviceHandler *ServiceHandler) Clear(c *gin.Context) {
	err := serviceHandler.ServiceUseCase.Clear()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Status(http.StatusOK)
}

func (serviceHandler *ServiceHandler) GetStatus(c *gin.Context) {
	status, err := serviceHandler.ServiceUseCase.GetStatus()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	statusJSON, err := status.MarshalJSON()
	if err != nil {
		c.Data(errors.PrepareErrorResponse(err))
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", statusJSON)
}
