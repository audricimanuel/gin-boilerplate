package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gin-boilerplate/src/config"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

type (
	GoMiddleware interface {
		LogRequest() gin.HandlerFunc
		RecoverPanic() gin.HandlerFunc
		BasicAuth(username, password string) gin.HandlerFunc
	}

	GoMiddlewareImpl struct {
		Config config.Config
	}
)

const (
	ParamQueryPage    = "page"
	ParamQueryLimit   = "limit"
	ParamQueryOffset  = "offset"
	ParamQueryKeyword = "keyword"
)

func InitMiddleware(cfg config.Config) GoMiddleware {
	return &GoMiddlewareImpl{
		Config: cfg,
	}
}

func (m *GoMiddlewareImpl) LogRequest() gin.HandlerFunc {
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

func (m *GoMiddlewareImpl) RecoverPanic() gin.HandlerFunc {
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

func (m *GoMiddlewareImpl) BasicAuth(username, password string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		if auth == "" {
			ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Decode the Base64 encoded auth string
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		userPass := strings.SplitN(string(decoded), ":", 2)
		if len(userPass) != 2 || userPass[0] != username || userPass[1] != password {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		ctx.Next()
	}
}
