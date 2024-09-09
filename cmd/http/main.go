package main

import (
	"context"
	"fmt"
	"gin-boilerplate/docs"
	"gin-boilerplate/src/config"
	"gin-boilerplate/src/database"
	"gin-boilerplate/src/internals/controller"
	"gin-boilerplate/src/internals/repository"
	"gin-boilerplate/src/internals/service"
	httpServer "gin-boilerplate/src/server/http"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func setSwaggerInfo() {
	docs.SwaggerInfo.Title = "Gin Boilerplate Example"
	docs.SwaggerInfo.Description = "Example boilerplate"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
}

func main() {
	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	// initialize mongodb connection
	databaseCollection := database.NewDatabaseCollection(cfg)

	// repositories
	exampleRepo := repository.NewExampleRepository(databaseCollection)

	// services
	exampleService := service.NewExampleService(exampleRepo)

	// controllers
	exampleController := controller.NewExampleController(exampleService)

	// set swagger info
	setSwaggerInfo()

	// registering router
	router := httpServer.RegisterRouter(
		cfg,
		exampleController,
		// register controllers in here
	)

	// running server
	logrus.Println("[INFO] Loading server")
	runServer(cfg, router)
}

func runServer(cfg config.Config, route http.Handler) {
	// The HTTP Server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host.Address, cfg.Host.Port),
		WriteTimeout: time.Second * time.Duration(cfg.Host.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(cfg.Host.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.Host.IdleTimeout),
		Handler:      route,
	}

	// Run Server
	go func() {
		logrus.Printf("⇨ http server started on %s\n", server.Addr)
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	logrus.Println("received shutdown signal. Trying to shutdown gracefully", sig)

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop Server
	logrus.Println("Stopping HTTP Server")
	server.SetKeepAlivesEnabled(false)
	err := server.Shutdown(ctx)
	if err != nil {
		logrus.Fatal("Failure while shutting down gracefully, errApp: ", err)
	}

	logrus.Println("Shutdown gracefully completed")
}
