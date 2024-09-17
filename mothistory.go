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

var errorMessages = map[int]string{
	400: "Bad Request - The data provided in the request is invalid. Please check your parameter and ensure that any date provided is no more than 5 weeks prior to today's date.",
	403: "Unauthorised - The API key is invalid. Please provide a valid API key.",
	404: "Resource Not Found - Vehicle with the provided parameters was not found or its test records are not valid. Please check your parameters.",
	429: "Too Man Requests - You have exceed your rate limit or quota.  Please wait before making additional requests.",
	500: "Internal Server Error - An unexpected issue occurred on server. Please contact support and provide the request ID returned.",
	503: "Service Unavailable - The service is unavailable. Please try again later.",
	504: "Gateway Timeout - The service did not respond within the time limit. Please try again later.",
}

func (c *Client) get(endpoint string, queryParams map[string]string) (*simplejson.Json, error) {
	response, err := c.restClient.R().
		SetHeader("x-api-key", c.apiKey).
		SetQueryParams(queryParams).
		Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to make request to API endpoint: %v", err)
	}

	errMsg, found := errorMessages[response.StatusCode()]

	if found {
		return nil, fmt.Errorf(errMsg, response.StatusCode())
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
