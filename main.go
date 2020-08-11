package main

import (
	"log"

	"github.com/egroj97/Golang-Demo/app"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//@TODO: Improve info fetching on the DB package.
//@TODO: Create pagination.
//@TODO: Improve logging.

func main() {
	a, err := app.New()
	if err != nil {
		panic(err)
	}
	log.Fatal(a.Run())
}
