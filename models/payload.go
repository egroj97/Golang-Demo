package models

// Payload will store all of the information obtained from the request.
type Payload struct {
	Model
	ConformsTo  string      `json:"conformsTo,omitempty"`
	DescribedBy string      `json:"describedBy,omitempty"`
	Context     string      `json:"@context,omitempty"`
	ElemType    string      `json:"@type,omitempty"`
	Dataset     []DataEntry `json:"dataset,omitempty"`
}

// DataEntry will store all the information on each entry of the  dataset field.
type DataEntry struct {
	Model
	PayloadID     uint           `json:"-"`
	ElemType      string         `json:"@type,omitempty"`
	Title         string         `json:"title,omitempty"`
	Description   string         `json:"description,omitempty"`
	Modified      string         `json:"modified,omitempty"`
	AccessLevel   string         `json:"accessLevel,omitempty"`
	Identifier    string         `json:"identifier,omitempty"`
	License       string         `json:"license,omitempty"`
	Publisher     Publisher      `json:"publisher,omitempty"`
	ContactPoint  ContactPoint   `json:"contactPoint,omitempty"`
	Distributions []Distribution `json:"distribution,omitempty"`
	Keywords      string         `json:"keywords,omitempty"`
	BureauCodes   string         `json:"bureauCodes,omitempty"`
	ProgramCodes  string         `json:"programCodes,omitempty"`
}

// Publisher will store all the information on the publisher field.
type Publisher struct {
	Model
	DataEntryID uint   `json:"-"`
	ElemType    string `json:"@type,omitempty"`
	Name        string `json:"name,omitempty"`
}

// ContactPoint will store all the information on the contactPoint field.
type ContactPoint struct {
	Model
	DataEntryID uint   `json:"-"`
	ElemType    string `json:"@type,omitempty"`
	Fn          string `json:"fn,omitempty"`
	HasEmail    string `json:"hasEmail,omitempty"`
}

// Distribution will store all the information on the distribution field.
type Distribution struct {
	Model
	DataEntryID uint   `json:"-"`
	ElemType    string `json:"@type,omitempty"`
	MediaType   string `json:"media_type,omitempty"`
	Format      string `json:"format,omitempty"`
	Title       string `json:"title,omitempty"`
	ConformsTo  string `json:"conforms_to,omitempty"`
	DownloadURL string `json:"download_url,omitempty"`
	AccessURL   string `json:"access_url,omitempty"`
}
