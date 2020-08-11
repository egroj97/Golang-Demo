package app

import (
	"net/http"

	"github.com/egroj97/Golang-Demo/models"

	"github.com/labstack/echo/v4"
)

// PageInfoHandler handles the rendering and retrieval of the information to be
// shown.
func (app *App) PageInfoHandler(c echo.Context) error {
	info := make([]models.Payload, 0)
	//Crashes on too many values, must change
	err := app.db.Set("gorm:auto_preload", true).Find(&info).Error
	if err != nil {
		return err
	}

	page := models.PageData{Payloads: info}

	return c.Render(http.StatusOK, "info", page)
}

// RestDataHandler handles the parsing and retrieval of the information to be
// sent if the JSON response
func (app *App) RestDataHandler(c echo.Context) error {
	info := make([]models.Payload, 0)
	//Crashes on too many values, must change
	err := app.db.Set("gorm:auto_preload", true).Find(&info).Error
	if err != nil {
		return err
	}

	page := models.PageData{Payloads: info}

	return c.JSON(http.StatusOK, page)
}
