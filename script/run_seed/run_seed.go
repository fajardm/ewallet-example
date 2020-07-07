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
	dbName := viper.GetString("DATABASE.NAME")
	conn, err := sql.Open(`mysql`, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error connecting database"))
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error ping database"))
	}
	files, err := ioutil.ReadDir("./database/seeds")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error read seeds directory"))
	}
	for _, file := range files {
		f, err := ioutil.ReadFile("./database/seeds/" + file.Name())
		fmt.Println(file.Name())
		if err != nil {
			log.Fatal(errors.Wrap(err, "Fatal error read seed file"))
		}
		scripts := strings.Split(string(f), ";")
		for _, script := range scripts {
			script := strings.TrimSpace(script)
			if script != "" {
				if _, err := conn.Exec(script); err != nil {
					log.Fatal(errors.Wrap(err, "Fatal error exec seed file"))
				}
			}
		}
	}
	db := &database.MySQL{DB: conn}
	if err := db.Close(); err != nil {
		log.Fatal(errors.Wrap(err, "Fatal error close database"))
	}
}
