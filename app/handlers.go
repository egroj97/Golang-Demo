package app

import (
	"net/http"
	"strings"

	"github.com/egroj97/Golang-Demo/models"

	"github.com/labstack/echo/v4"
)

// PageInfoHandler handles the rendering and retrieval of the information to be
// shown.
func (app *App) PageInfoHandler(c echo.Context) error {
	app.server.Logger.Info("Reading from DB")
	info, err := app.fetchAllPayloadRecords()
	if err != nil {
		return err
	}

	page := models.PageData{Payloads: info}

	app.server.Logger.Info("Rendering template")
	return c.Render(http.StatusOK, "info", page)
}

// RestDataHandler handles the parsing and retrieval of the information to be
// sent if the JSON response
func (app *App) RestDataHandler(c echo.Context) error {
	app.server.Logger.Info("Reading from DB")
	info, err := app.fetchAllPayloadRecords()
	if err != nil {
		return err
	}

	resp := make([]map[string]interface{}, len(info))
	for i, payload := range info {
		resp[i] = make(map[string]interface{})
		resp[i]["conformsTo"] = payload.ConformsTo
		resp[i]["describedBy"] = payload.DescribedBy
		resp[i]["@context"] = payload.Context
		resp[i]["@type"] = payload.ElemType
		dataset := make([]models.DataEntryJSON, len(payload.Dataset))
		for j, dataEntry := range payload.Dataset {
			temp := models.DataEntryJSON{DataEntry: &dataEntry}
			temp.Keywords = strings.Split(dataEntry.Keywords, ", ")
			temp.ProgramCodes = strings.Split(dataEntry.ProgramCodes, ", ")
			temp.BureauCodes = strings.Split(dataEntry.BureauCodes, ", ")

			dataset[j] = temp
		}
		resp[i]["dataset"] = dataset
	}

	app.server.Logger.Info("Sending JSON response")
	return c.JSON(http.StatusOK, resp)
}
