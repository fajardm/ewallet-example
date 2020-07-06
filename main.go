package main

import (
	"database/sql"
	"fmt"
	_userHttp "github.com/fajardm/ewallet-example/app/user/http"
	_userRepository "github.com/fajardm/ewallet-example/app/user/repository/mysql"
	_userUsecase "github.com/fajardm/ewallet-example/app/user/usecase"
	"github.com/fajardm/ewallet-example/bootstrap"
	_ "github.com/go-sql-driver/mysql"
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

func prepareDatabase() *sql.DB {
	dbUser := viper.GetString("DATABASE.USER")
	dbPassword := viper.GetString("DATABASE.PASSWORD")
	dbHost := viper.GetString("DATABASE.HOST")
	dbPort := viper.GetString("DATABASE.PORT")
	dbName := viper.GetString("DATABASE.NAME")
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	conn, err := sql.Open(`mysql`, connStr)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error connecting database"))
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error ping database"))
	}
	return conn
}

func main() {
	prepareConfig()
	contextTimeout := viper.GetDuration("CONTEXT_TIMEOUT")

	conn := prepareDatabase()
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "Fatal error close database"))
		}
	}()

	app := bootstrap.New(viper.GetString("APP_NAME"), viper.GetString("APP_OWNER"))
	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("Ok!")
	})

	// Register user handler
	userRepository := _userRepository.NewUserRepository(conn)
	userUsecase := _userUsecase.NewUserUsecase(userRepository, contextTimeout)
	_userHttp.NewUserHandler(app, userUsecase)

	if err := app.Listen(viper.GetInt("APP_PORT")); err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error listen port"))
	}
}
