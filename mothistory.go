package mothistory

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	// "strings"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"
)

var BaseURL string = "https://history.mot.api.gov.uk/v1/trade/vehicles"

const (
	TokenURL   = "https://login.microsoftonline.com/a455b827-244f-4c97-b5b4-ce5d13b4d00c/oauth2/v2.0/token"
	ScopeURL   = "https://tapi.dvsa.gov.uk/.default"
	RPSLimit   = 15
	BurstLimit = 10
	DailyQuota = 500000
)

type Client struct {
	httpClient  *http.Client
	apiKey      string
	rateLimiter rate.Limiter
	dayLimiter  rate.Limiter
}

type ClientConfig struct {
	ClientID     string
	ClientSecret string
	APIKey       string
}

func NewClient(config ClientConfig, customHTTPClient *http.Client) *Client {
	oauthConfig := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     TokenURL,
		Scopes:       []string{ScopeURL},
	}

	dailyLimit := rate.Limit(float64(DailyQuota) / float64(23*60*60)) // Daily Limit = DailyQuota / number of seconds in a day

	// Use custom HTTP client if provided.
	httpClient := customHTTPClient
	if httpClient == nil {
		httpClient = oauthConfig.Client(context.Background())
	}

	return &Client{
		httpClient:  httpClient,
		apiKey:      config.APIKey,
		rateLimiter: *rate.NewLimiter(RPSLimit, BurstLimit),
		dayLimiter:  *rate.NewLimiter(dailyLimit, DailyQuota),
	}
}

var errorMessages = map[int]string{
	400: "Bad Request - The format of the request is incorrect",
	401: "Unauthorized - Authentication credentials are missing or invalid",
	403: "Forbidden - The request is not allowed",
	404: "Not Found - The requested data is not found",
	405: "Method Not Allowed - The HTTP method is not supported for this endpoint",
	406: "Not Acceptable - The requested media type is not supported",
	409: "Conflict - The request could not be completed due to a conflict with the current state of the target resource",
	412: "Precondition Failed - Could not complete request because a constraint was not met",
	415: "Unsupported Media Type - The media type of the request is not supported",
	422: "Unprocessable Entity - The request was well-formed but contains semantic errors",
	429: "Too Many Requests - The user has sent too many requests in a given amount of time",
	500: "Internal Server Error - An unexpected error has occurred",
	502: "Bad Gateway - The server received an invalid response from an upstream server",
	503: "Service Unavailable - The server is currently unable to handle the request",
	504: "Gateway Timeout - The upstream server failed to send a request in the time allowed by the server",
}

func doRequest[T any](c *Client, method, endpoint string, queryParams url.Values)(*T, error) {
	limiterCtx := context.Background()
	if err := c.dayLimiter.Wait(limiterCtx); err != nil {
		return nil, fmt.Errorf("daily quota exceeded: %v", err)
	}
	if err := c.rateLimiter.Wait(limiterCtx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %v", err)
	}

	url := fmt.Sprintf("%s%s", BaseURL, endpoint)
	if len(queryParams) > 0 {
		url = fmt.Sprintf("%s?%s", url, queryParams.Encode())
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if errMsg, found := errorMessages[resp.StatusCode]; found {
		return nil, fmt.Errorf("%s", errMsg)
	}

	var structuredResponse T
	err = json.Unmarshal(body, &structuredResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %v", err)
	}

	return &structuredResponse, nil
}

func (c *Client) GetByRegistration(registration string) (*VehicleDetails, error) {
	endpoint := fmt.Sprintf("/registration/%s", url.PathEscape(registration))
	return doRequest[VehicleDetails](c, http.MethodGet, endpoint, nil)
}

func (c *Client) GetByVIN(vin string) (*VehicleDetails, error) {
	endpoint := fmt.Sprintf("/vin/%s", url.PathEscape(vin))
	return doRequest[VehicleDetails](c, http.MethodGet, endpoint, nil)
}

func (c *Client) GetBulkDownload() (*BulkDownloadResponse, error) {
	return doRequest[BulkDownloadResponse](c, http.MethodGet, "/bulk-download", nil)
}

// func (c *Client) RenewCredentials(apiKeyValue, email string) (json.RawMessage, error) {
// 	payload := url.Values{}
// 	payload.Set("awsApiKeyValue", apiKeyValue)
// 	payload.Set("email", email)

// 	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/credentials", BaseURL), strings.NewReader(payload.Encode()))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %v", err)
// 	}

// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Set("X-API-Key", c.apiKey)

// 	resp, err := c.httpClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to make request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response body: %v", err)
// 	}

// 	if errMsg, found := errorMessages[resp.StatusCode]; found {
// 		return nil, fmt.Errorf("%s", errMsg)
// 	}

// 	return body, nil
// }
