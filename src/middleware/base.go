package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-boilerplate/src/config"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"runtime/debug"
)

type (
	GoMiddleware struct {
		Config config.Config
	}
)

const (
	ParamQueryPage    = "page"
	ParamQueryLimit   = "limit"
	ParamQueryOffset  = "offset"
	ParamQueryKeyword = "keyword"
)

func InitMiddleware(cfg config.Config) *GoMiddleware {
	return &GoMiddleware{
		Config: cfg,
	}
}

func (m *GoMiddleware) LogRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestLog := MapLogRequest(ctx)
		fmt.Println(requestLog)

		ctx.Next()
	}
}

// MapLogRequest for map log request
func MapLogRequest(ctx *gin.Context) string {
	if ctx.GetHeader("Content-Type") == "application/json" {
		// Read the content
		var bodyBytes []byte
		if ctx.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(ctx.Request.Body)
		}
		// Use the content
		var req interface{}
		json.Unmarshal(bodyBytes, &req)
		bodyBytes, _ = json.Marshal(req)

		// Restore the io.ReadCloser to its original state
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return fmt.Sprintf("[IN_REQUEST: [%s] %s]", ctx.Request.Method, ctx.Request.URL.String())
}

func ServerError(ctx *gin.Context, err error, code int) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Output(2, trace)
	ctx.Set("Content-Type", "application/json")

	ctx.JSON(code, http.StatusText(code))
}

func (m *GoMiddleware) RecoverPanic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Set("Connection", "close")
				ServerError(ctx, fmt.Errorf("%s", err), http.StatusInternalServerError)
			}
		}()

		ctx.Next()
	}
}
