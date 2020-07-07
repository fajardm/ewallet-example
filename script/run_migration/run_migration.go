package main

import (
	"database/sql"
	"fmt"
	"github.com/fajardm/ewallet-example/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error config file"))
	}

	dbUser := viper.GetString("DATABASE.USER")
	dbPassword := viper.GetString("DATABASE.PASSWORD")
	dbHost := viper.GetString("DATABASE.HOST")
	dbPort := viper.GetString("DATABASE.PORT")
	//dbName := viper.GetString("DATABASE.NAME")
	conn, err := sql.Open(`mysql`, fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true", dbUser, dbPassword, dbHost, dbPort))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error connecting database"))
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error ping database"))
	}
	files, err := ioutil.ReadDir("./database/migrations")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error read migrations directory"))
	}
	for _, file := range files {
		f, err := ioutil.ReadFile("./database/migrations/" + file.Name())
		if err != nil {
			log.Fatal(errors.Wrap(err, "Fatal error read migration file"))
		}
		scripts := strings.Split(string(f), ";")
		for _, script := range scripts {
			script := strings.TrimSpace(script)
			if script != "" {
				if _, err := conn.Exec(script); err != nil {
					log.Fatal(errors.Wrap(err, "Fatal error exec migration file"))
				}
			}
		}
	}
	db := &database.MySQL{DB: conn}
	if err := db.Close(); err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error close database"))
	}
}
