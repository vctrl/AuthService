package main

import (
	"log"

	"github.com/spf13/viper"
	"github.com/vctrl/authService/config"
	"github.com/vctrl/authService/server"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	app, err := server.NewApp()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	if err := app.Run(viper.GetString("port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
