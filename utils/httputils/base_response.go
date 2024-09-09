package httputils

import (
	"github.com/audricimanuel/errorutils"
	"github.com/gin-gonic/gin"
	"math"
)

type (
	BaseMeta struct {
		Page      int `json:"page"`
		Limit     int `json:"limit"`
		TotalData int `json:"total_data"`
		TotalPage int `json:"total_page"`
	}

	// BaseResponse is the base response
	BaseResponse struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data"`
		Error  *string     `json:"error_message"`
		Meta   *BaseMeta   `json:"meta,omitempty"`
	}
)

// SetHttpResponse map response
func SetHttpResponse(ctx *gin.Context, data interface{}, err errorutils.HttpError, meta *BaseMeta) {
	var errMsg *string

	statusCode, message := errorutils.GetStatusCode(err)
	if message != errorutils.SUCCESS {
		errMsg = &message
	}

	ctx.JSON(statusCode, BaseResponse{
		Status: statusCode,
		Data:   data,
		Error:  errMsg,
		Meta:   meta,
	})
}

func SetBaseMeta(page int, limit int, totalData int) BaseMeta {
	totalPage := float64(totalData) / float64(limit)
	return BaseMeta{
		Page:      page,
		Limit:     limit,
		TotalData: totalData,
		TotalPage: int(math.Ceil(totalPage)),
	}
}
