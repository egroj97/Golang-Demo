package main

import (
	"github.com/egroj97/Golang-Demo/app"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/gommon/log"
)

//@TODO: Create pagination.

func main() {
	a, err := app.New()
	if err != nil {
		panic(err)
	}
	log.Fatal(a.Run())
}
