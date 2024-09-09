package controller

import (
	"gin-boilerplate/src/internals/service"
	"gin-boilerplate/utils/httputils"
	"github.com/gin-gonic/gin"
)

type (
	ExampleController interface {
		GetExample(ctx *gin.Context)
	}

	ExampleControllerImpl struct {
		exampleService service.ExampleService
	}
)

func NewExampleController(e service.ExampleService) ExampleController {
	return &ExampleControllerImpl{
		exampleService: e,
	}
}

// @Tags			Example
// @Summary		Example API
// @Description	"Just an example"
// @Accept			json
// @Produce		json
// @Success		200	{object}	httputils.BaseResponse
// @Router			/example [get]
func (e *ExampleControllerImpl) GetExample(ctx *gin.Context) {
	data := e.exampleService.GetExample(ctx)
	httputils.SetHttpResponse(ctx, data, nil, nil)
}
