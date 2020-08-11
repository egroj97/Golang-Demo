package app

import (
	"html/template"
	"time"

	"github.com/egroj97/Golang-Demo/models"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// URL from which the information will be retrieved.
const URL = "http://www.energy.gov/data.json"

// App contains all of the elements that the application will need to run
// successfully.
type App struct {
	db     *gorm.DB
	server *echo.Echo
	client *resty.Client
}

// Run executes our main application.
func (app *App) Run() error {
	// A ticker is created to fetch updated data every 24 hours.
	fetchTicker := time.NewTicker(24 * time.Hour)
	defer fetchTicker.Stop()

	done := make(chan bool)

	// A first payload fetch is made if there aren't any payloads on the DB.
	payload := models.Payload{}
	if app.db.First(&payload).RecordNotFound() {
		err := app.fetchPayload()
		if err != nil {
			return err
		}
	}

	// Concurrently a new payload will be requested every 24 hours as per the
	// ticker created at the beginning.
	go func() {
		for true {
			select {
			case <-done:
				return
			case <-fetchTicker.C:
				// This call to fetchPayload will be done on "best effort"
				// logging errors encountered instead of failure.
				err := app.fetchPayload()
				if err != nil {
					err := app.fetchPayload()
					app.server.Logger.Errorf("Could not fetch data: %v", err)
				}
			}
		}
	}()

	// Separately the server will be run on its own goroutine .
	app.setupServer()
	return app.server.Start(":1323")
}

// fetchPayload will retrieve the data from the selected URL.
func (app *App) fetchPayload() error {
	resp, err := app.client.R().Get(URL)
	if err != nil {
		return err
	}

	payload, err := parsePayload(resp.Body())
	if err != nil {
		return err
	}

	err = app.db.Create(&payload).Error
	if err != nil {
		return err
	}

	return nil
}

// setupServer will configure all the endpoints needed alongside the template
// engine.
func (app *App) setupServer() {
	t := &models.Template{
		Templates: template.Must(template.ParseGlob("static/*.html")),
	}

	app.server.Renderer = t
	app.server.Use(middleware.Logger())
	app.server.GET("/", app.PageInfoHandler)
	app.server.GET("/data", app.RestDataHandler)
}
