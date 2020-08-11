package app

// @TODO: Better refactoring of this functionality

import (
	json "github.com/buger/jsonparser"
	"github.com/egroj97/Golang-Demo/models"
)

// parsePayload parses the information inside of the request into a Payload
// structure.
func parsePayload(data []byte) (models.Payload, error) {
	var payload models.Payload
	var err error

	payload.ElemType, err = json.GetString(data, "@type")
	if err != nil {
		return models.Payload{}, err
	}

	payload.ConformsTo, err = json.GetString(data, "conformsTo")
	if err != nil {
		return models.Payload{}, err
	}

	payload.DescribedBy, err = json.GetString(data, "describedBy")
	if err != nil {
		return models.Payload{}, err
	}

	payload.Context, err = json.GetString(data, "@context")
	if err != nil {
		return models.Payload{}, err
	}

	entries, dataType, _, err := json.Get(data, "dataset")
	if err != nil {
		panic(err)
	}

	if dataType == json.Array {
		_, err = json.ArrayEach(
			entries,
			func(value []byte, dataType json.ValueType, offset int, err error) {
				dataEntry := &models.DataEntry{}
				err = parseEntry(dataEntry, value)
				if err != nil {
					return
				}
				if dataEntry.ElemType != "" {
					payload.Dataset = append(payload.Dataset, *dataEntry)
				}
			})
		if err != nil {
			return models.Payload{}, err
		}
	}

	return payload, nil
}

// parseEntry parses each entry of the dataset field inside of a DataEntry
// structure.
func parseEntry(result *models.DataEntry, data []byte) (err error) {
	result.ElemType, err = json.GetString(data, "@type")
	if err != nil {
		return err
	}

	result.Title, err = json.GetString(data, "title")
	if err != nil {
		return err
	}

	result.Description, err = json.GetString(data, "description")
	if err != nil {
		return err
	}

	result.Modified, err = json.GetString(data, "modified")
	if err != nil {
		return err
	}

	result.AccessLevel, err = json.GetString(data, "accessLevel")
	if err != nil {
		return err
	}

	result.Identifier, err = json.GetString(data, "identifier")
	if err != nil {
		return err
	}

	result.License, err = json.GetString(data, "license")
	if err != nil {
		return err
	}

	result.Publisher.ElemType, err = json.GetString(data, "publisher", "@type")
	if err != nil {
		return err
	}

	result.Publisher.Name, err = json.GetString(data, "publisher", "name")
	if err != nil {
		return err
	}

	result.ContactPoint.ElemType, err = json.GetString(data, "contactPoint", "@type")
	if err != nil {
		return err
	}

	result.ContactPoint.Fn, err = json.GetString(data, "contactPoint", "fn")
	if err != nil {
		return err
	}

	result.ContactPoint.HasEmail, err = json.GetString(data, "contactPoint", "hasEmail")
	if err != nil {
		return err
	}

	distribution, dataType, _, err := json.Get(data, "distribution")
	if err != nil {
		return err
	}

	if dataType == json.Array {
		result.Distributions = make([]models.Distribution, 0)
		_, err = json.ArrayEach(
			distribution,
			func(value []byte, dataType json.ValueType, offset int, err error) {
				dist := models.Distribution{}
				dist.ElemType, err = json.GetString(value, "@type")
				dist.MediaType, err = json.GetString(value, "mediaType")
				dist.Format, err = json.GetString(value, "format")
				dist.Title, err = json.GetString(value, "title")
				dist.ConformsTo, err = json.GetString(value, "conformsTo")
				dist.DownloadURL, err = json.GetString(value, "downloadURL")
				dist.AccessURL, err = json.GetString(value, "accessURL")
				result.Distributions = append(result.Distributions, dist)
			})
		if err != nil {
			return err
		}
	}

	keyword, dataType, _, err := json.Get(data, "keyword")
	if dataType == json.Array {
		_, err = json.ArrayEach(
			keyword,
			func(value []byte, dataType json.ValueType, offset int, err error) {
				if dataType == json.String {
					if result.Keywords != "" {
						result.Keywords += ", "
					}
					result.Keywords += string(value)
				}
			})
		if err != nil {
			return err
		}
	}

	bureauCode, dataType, _, err := json.Get(data, "bureauCode")
	if dataType == json.Array {
		_, err = json.ArrayEach(
			bureauCode,
			func(value []byte, dataType json.ValueType, offset int, err error) {
				if dataType == json.String {
					if result.BureauCodes != "" {
						result.BureauCodes += ", "
					}
					result.BureauCodes += string(value)
				}
			})
		if err != nil {
			return err
		}
	}

	programCode, dataType, _, err := json.Get(data, "programCode")
	if dataType == json.Array {
		_, err = json.ArrayEach(
			programCode,
			func(value []byte, dataType json.ValueType, offset int, err error) {
				if dataType == json.String {
					if result.ProgramCodes != "" {
						result.ProgramCodes += ", "
					}
					result.ProgramCodes += string(value)
				}
			})
		if err != nil {
			return err
		}
	}

	return nil
}
