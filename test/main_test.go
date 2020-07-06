package main

import (
	"database/sql"
	"fmt"
	_userHttp "github.com/fajardm/ewallet-example/app/user/http"
	_userRepository "github.com/fajardm/ewallet-example/app/user/repository/mysql"
	_userUsecase "github.com/fajardm/ewallet-example/app/user/usecase"
	"github.com/fajardm/ewallet-example/bootstrap"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
	"testing"
)

var app *bootstrap.Bootstrap

func TestMain(m *testing.M) {
	viper.SetConfigFile("../config.test.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error config file"))
	}

	contextTimeout := viper.GetDuration("CONTEXT_TIMEOUT")

	dbUser := viper.GetString("DATABASE.USER")
	dbPassword := viper.GetString("DATABASE.PASSWORD")
	dbHost := viper.GetString("DATABASE.HOST")
	dbPort := viper.GetString("DATABASE.PORT")
	dbName := viper.GetString("DATABASE.NAME")
	conn, err := sql.Open(`mysql`, fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true", dbUser, dbPassword, dbHost, dbPort))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error connecting database"))
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error ping database"))
	}
	if _, err := conn.Exec("DROP DATABASE IF EXISTS " + dbName); err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error drop database"))
	}
	files, err := ioutil.ReadDir("../database/migrations")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error read migrations directory"))
	}
	for _, file := range files {
		f, err := ioutil.ReadFile("../database/migrations/" + file.Name())
		if err != nil {
			log.Fatal(errors.Wrap(err, "Fatal error read migration file"))
		}
		scripts := strings.Split(strings.Replace(string(f), "ewallet", dbName, -1), ";")
		for _, script := range scripts {
			script := strings.TrimSpace(script)
			if script != "" {
				if _, err := conn.Exec(script); err != nil {
					log.Fatal(errors.Wrap(err, "Fatal error exec migration file"))
				}
			}
		}
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "Fatal error close database"))
		}
	}()

	app = bootstrap.New(viper.GetString("APP_NAME"), viper.GetString("APP_OWNER"))

	// Register user handler
	userRepository := _userRepository.NewUserRepository(conn)
	userUsecase := _userUsecase.NewUserUsecase(userRepository, contextTimeout)
	_userHttp.NewUserHandler(app, userUsecase)

	m.Run()
}
