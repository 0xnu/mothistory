package mothistory

import (
	"fmt"

	"github.com/bitly/go-simplejson"
	"github.com/go-resty/resty/v2"
)

const BaseURL = "https://beta.check-mot.service.gov.uk/trade/vehicles/mot-tests"

type Client struct {
	apiKey     string
	restClient *resty.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		restClient: resty.New(),
	}
}

func (c *Client) get(endpoint string, queryParams map[string]string) (*simplejson.Json, error) {
	response, err := c.restClient.R().
		SetHeader("x-api-key", c.apiKey).
		SetQueryParams(queryParams).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to make request to API endpoint: %v", err)
	}

	if response.StatusCode() != 200 {
		return nil, fmt.Errorf("API returned non-200 status code: %d", response.StatusCode())
	}

	json, err := simplejson.NewJson(response.Body())

	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return json, nil
}

func (c *Client) GetByRegistration(registration string) (*simplejson.Json, error) {
	queryParams := map[string]string{
		"registration": registration,
	}
	return c.get(BaseURL, queryParams)
}

func (c *Client) GetByPage(page int) (*simplejson.Json, error) {
	queryParams := map[string]string{
		"page": fmt.Sprintf("%d", page),
	}
	return c.get(BaseURL, queryParams)
}

func (c *Client) GetByDateAndPage(date string, page int) (*simplejson.Json, error) {
	queryParams := map[string]string{
		"date": date,
		"page": fmt.Sprintf("%d", page),
	}
	return c.get(BaseURL, queryParams)
}

func (c *Client) GetByVehicleID(vehicleId string) (*simplejson.Json, error) {
	queryParams := map[string]string{
		"vehicleId": vehicleId,
	}
	return c.get(BaseURL, queryParams)
}
