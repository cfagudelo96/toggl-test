// Package app contains the web application server
package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/cfagudelo96/toggle-test/deck/handler"
	"github.com/cfagudelo96/toggle-test/deck/repository"
	"github.com/cfagudelo96/toggle-test/deck/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// App represents the web application.
type App struct {
	Server *echo.Echo
}

// NewApp creates a new app and leaves it ready for execution.
func NewApp() *App {
	a := &App{}
	a.setupServer()

	return a
}

func (a *App) setupServer() {
	a.Server = echo.New()
	a.Server.Use(
		middleware.Recover(),
		middleware.RequestID(),
		middleware.Logger(),
	)
	a.setupRoutes()
}

func (a *App) setupRoutes() {
	dr := repository.NewInMemoryDeckRepository()
	ds := service.NewDeckService(dr)
	dh := handler.NewDeckEchoHandler(ds)
	apiGroup := a.Server.Group("/v1/decks")
	apiGroup.POST("", dh.HandleCreateDeck)
	apiGroup.GET("/:uuid", dh.HandleOpenDeck)
	apiGroup.POST("/:uuid/draw", dh.HandleDrawCars)
}

// StartApp initializes the server.
func (a *App) StartApp() {
	go a.startServer()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	a.stopApp()
}

func (a *App) startServer() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	if err := a.Server.Start(fmt.Sprintf(":%s", port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error starting the server: %v", err)
	}
}

func (a *App) stopApp() {
	ctx := context.Background()
	if err := a.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down the server: %v", err)
	}
}
