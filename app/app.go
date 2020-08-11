package app

import (
	"github.com/egroj97/Golang-Demo/models"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
)

// New configure and returns a pointer to an App struct with all of their
// dependencies
func New() (*App, error) {
	var err error
	app := &App{}

	app.client = resty.New()
	app.server = echo.New()

	app.db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	err = app.db.AutoMigrate(&models.Distribution{}, &models.Publisher{},
		&models.ContactPoint{}, &models.DataEntry{}, &models.Payload{}).Error
	if err != nil {
		return nil, err
	}

	return app, nil
}
