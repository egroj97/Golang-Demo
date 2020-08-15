package app

import (
	"html/template"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/egroj97/Golang-Demo/models"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// App contains all of the elements that the application will need to run
// successfully.
type App struct {
	db     *gorm.DB
	server *echo.Echo
	client *resty.Client
}

// New configure and returns a pointer to an App struct with all of its
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

	logger := log.New("")

	logger.SetHeader("[${time_rfc3339}] ${level} | ${short_file} | line=${line} |")

	app.client.SetLogger(logger)
	app.server.Logger = logger

	app.server.Renderer = &models.Template{
		Templates: template.Must(template.ParseGlob("static/*.html")),
	}

	app.server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method} | uri=${uri} | status=${status} |" +
			" latency_human=${latency_human} | bytes_in=${bytes_in} | " +
			"bytes_out=${bytes_out}\n",
	}))

	err = app.db.AutoMigrate(&models.Distribution{}, &models.Publisher{},
		&models.ContactPoint{}, &models.DataEntry{}, &models.Payload{}).Error
	if err != nil {
		return nil, err
	}

	return app, nil
}

// Run executes our main application.
func (app *App) Run() error {
	// A ticker is created to fetch updated data every 24 hours.
	fetchTicker := time.NewTicker(24 * time.Hour)
	defer fetchTicker.Stop()

	done := make(chan bool)

	// A first payload fetch is made if there aren't any payloads on the DB.
	payload := models.Payload{}
	if app.db.Last(&payload).RecordNotFound() ||
		time.Now().Sub(payload.CreatedAt) >= 24*time.Hour {
		app.server.Logger.Info("Fetching payload")
		err := app.fetchPayload()
		if err != nil {
			return err
		}
	}

	// Concurrently a new payload will be requested every 24 hours as per the
	// ticker created at the beginning.
	app.server.Logger.Info("Launching fetcher routine")
	go func() {
		for true {
			select {
			case <-done:
				return
			case <-fetchTicker.C:
				// This call to fetchPayload will be done on "best effort"
				// logging errors encountered instead of failure.
				err := app.fetchPayload()
				app.server.Logger.Info("Fetching payload")
				if err != nil {
					app.server.Logger.Errorf("Could not fetch data: %v", err)
				}
			}
		}
	}()

	// Separately the server will be run on its own goroutine .
	app.setupEndpoints()
	return app.server.Start(":1323")
}
