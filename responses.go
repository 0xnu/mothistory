package mothistory

import (
	// "encoding/json"
)
// Response for {baseURL}/[registration|vin]/<registration|vin>
type VehicleDetails struct {
	Registration         string    `json:"registration"`
	Make                 string    `json:"make"`
	FirstUsedDate        string    `json:"firstUsedDate"`
	FuelType             string    `json:"fuelType"`
	PrimaryColour        string    `json:"primaryColour"`
	RegistrationDate     string    `json:"registrationDate"`
	ManufactureDate      string    `json:"manufactureDate"`
	EngineSize           string    `json:"engineSize"`
	HasOutstandingRecall string    `json:"hasOutstandingRecall"`
	MotTests             []MOTTest `json:"motTests"`
}

type MOTTest struct {
	CompletedDate      string   `json:"completedDate"`
	TestResult         string   `json:"testResult"`
	ExpiryDate         string   `json:"expiryDate"`
	OdometerValue      string   `json:"odometerValue"`
	OdometerUnit       string   `json:"odometerUnit"`
	OdometerResultType string   `json:"odometerResultType"`
	MotTestNumber      string   `json:"motTestNumber"`
	DataSource         string   `json:"dataSource"`
	Location           string   `json:"location,omitempty"` // Optional field
	Defects            []Defect `json:"defects,omitempty"`  // Optional field
}

type Defect struct {
	Text         string `json:"text"`
	TypeOfDefect string `json:"type"`
	Dangerous    bool   `json:"dangerous"`
}

// Response for {baseURL}/[bulk-download]
type BulkDownload struct {
	Bulk  []BulkDelta `json:"bulk"`
	Delta []BulkDelta `json:"delta"`
}

type BulkDelta struct {
	Filename      string `json:"filename"`
	DonwloadURL   string `json:"downloadUrl"`
	FileSize      int    `json:"fileSize"`
	FileCreatedOn string `json:"fileCreatedOn"`
}
