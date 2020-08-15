package app

import (
	"github.com/egroj97/Golang-Demo/models"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// URL from which the information will be retrieved.
const URL = "http://www.energy.gov/data.json"

// fetchPayload will retrieve the data from the selected URL.
func (app *App) fetchPayload() error {
	app.server.Logger.Info("Retrieving payload")
	resp, err := app.client.R().Get(URL)
	if err != nil {
		return err
	}

	app.server.Logger.Info("Parsing payload")
	payload, err := parsePayload(resp.Body())
	if err != nil {
		return err
	}

	app.server.Logger.Info("Storing to database")
	err = app.db.Create(&payload).Error
	if err != nil {
		return err
	}

	return nil
}

// setupEndpoints will configure all the endpoints needed alongside the template
// engine.
func (app *App) setupEndpoints() {
	app.server.GET("/", app.PageInfoHandler)
	app.server.GET("/data", app.RestDataHandler)
}

// fetchAllPayloadRecords loads all the records on DB without preloading them with
// GORM, said procedure fails on too many records.
func (app *App) fetchAllPayloadRecords() (records []models.Payload, err error) {
	err = app.db.Find(&records).Error
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(records); i++ {
		dataset := make([]models.DataEntry, 0)
		err := app.db.Where("payload_id = ?", records[i].ID).Find(&dataset).Error
		if err != nil {
			return nil, err
		}

		records[i].Dataset = dataset
		for j := 0; j < len(records[i].Dataset); j++ {
			publisher := models.Publisher{}
			err = app.db.Where(
				"data_entry_id = ?",
				records[i].Dataset[j].ID,
			).Find(&publisher).Error
			if err != nil {
				return nil, err
			}
			records[i].Dataset[j].Publisher = publisher

			contactPoint := models.ContactPoint{}
			err = app.db.Where(
				"data_entry_id = ?",
				records[i].Dataset[j].ID,
			).First(&contactPoint).Error
			if err != nil {
				return nil, err
			}
			records[i].Dataset[j].ContactPoint = contactPoint

			distributions := make([]models.Distribution, 0)
			err = app.db.Where(
				"data_entry_id = ?",
				records[i].Dataset[j].ID,
			).Find(&distributions).Error
			if err != nil {
				return nil, err
			}
			records[i].Dataset[j].Distributions = distributions
		}
	}

	return records, nil
}
