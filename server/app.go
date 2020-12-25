package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/vctrl/authService/delivery"
	"github.com/vctrl/authService/usecase"
)

type App struct {
	httpServer *http.Server
	authUC     *usecase.OAuthUseCase
}

func NewApp() *App {
	return &App{
		authUC: usecase.NewAuthUseCase(),
	}
}

func (a *App) Run(port string) error {
	// register auth endpoint
	delivery.RegisterHTTPEndpoints(a.authUC)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
