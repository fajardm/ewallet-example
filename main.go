package main

import (
	"github.com/fajardm/ewallet-example/bootstrap"
	"github.com/gofiber/fiber"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func prepareConfig() {
	file := os.Getenv("CONFIG")
	if file == "" {
		file = "config.yaml"
	}
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error config file"))
	}
}

func main() {
	prepareConfig()
	app := bootstrap.New(viper.GetString("APP_NAME"), viper.GetString("APP_OWNER"))
	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("Ok!")
	})
	if err := app.Listen(viper.GetInt("APP_PORT")); err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error listen port"))
	}
}
