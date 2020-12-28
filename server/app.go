package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/vctrl/authService/delivery"
	"github.com/vctrl/authService/middleware"
	"github.com/vctrl/authService/usecase"
)

type App struct {
	httpServer *http.Server
	authUC     usecase.UseCase
}

func NewApp() *App {
	configs := parseConfigs()

	return &App{
		authUC: usecase.NewAuthUseCase(configs),
	}
}

func (a *App) Run(port string) error {
	// install gin router
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.JSONAppErrorReporter(),
	)

	// register auth endpoint
	delivery.RegisterHTTPEndpoints(router, a.authUC)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
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

func parseConfigs() map[string]usecase.UseCaseConfig {
	var siteConfigs usecase.SiteConfigs
	err := viper.Unmarshal(&siteConfigs)

	configs := make(map[string]usecase.UseCaseConfig)
	for s, c := range siteConfigs {
		configs[s] = c.MapToOauth2Config(s)
	}

	if err != nil {
		fmt.Println("Error reading configs: ", err)
	}

	return configs
}
