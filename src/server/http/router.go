package http

import (
	"fmt"
	"gin-boilerplate/src/config"
	"gin-boilerplate/src/internals/controller"
	"gin-boilerplate/src/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func RegisterRouter(
	cfg config.Config,
	exampleController controller.ExampleController,
	// register new controllers here
) *gin.Engine {
	r := gin.Default()

	mid := middleware.InitMiddleware(cfg)

	setMiddlewareGlobal(mid, r)

	// Swagger
	r.Handle("GET", "/docs/*any", mid.BasicAuth(cfg.SwaggerUsername, cfg.SwaggerPassword), ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Handle("GET", "/docs", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})

	r.GET("/ping", func(ctx *gin.Context) {
		staticText := fmt.Sprintf("hello world: %s", cfg.Env)
		ctx.JSON(http.StatusOK, gin.H{"message": staticText})
	})

	r.GET("/example", exampleController.GetExample)

	return r
}

func setMiddlewareGlobal(mid middleware.GoMiddleware, r *gin.Engine) {
	// Logger
	r.Use(mid.LogRequest())

	// Cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Recovery
	r.Use(mid.RecoverPanic())
}
