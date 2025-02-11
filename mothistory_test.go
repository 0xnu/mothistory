package mothistory

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func createMockServer() *httptest.Server {
	handler := http.NewServeMux()

	handler.HandleFunc("/registration/ML58FOU", func(w http.ResponseWriter, r *http.Request) {
		mockResponse := `{"registration": "ML58FOU"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockResponse)
	})

	handler.HandleFunc("/vin/BNR32305366", func(w http.ResponseWriter, r *http.Request) {
		mockResponse := `{"registration": "ML58FOU"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockResponse)
	})

	handler.HandleFunc("/bulk-download", func(w http.ResponseWriter, r *http.Request) {
		mockResponse := `{"bulk": [], "delta": []}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockResponse)
	})

	handler.HandleFunc("/credentials", func(w http.ResponseWriter, r *http.Request) {
		mockResponse := `{"clientSecret": "your-new-client-secret-value"}`
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Header().Set("X-API-Key", "dummy-api-key")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockResponse)
	})

	mockServer := httptest.NewServer(handler)
	return mockServer
}

func TestGetByRegistration(t *testing.T) {
	mockServer := createMockServer()
	defer mockServer.Close()

	BaseURL = mockServer.URL
	client := createTestClient(mockServer)

	registration := "ML58FOU"
	data, err := client.GetByRegistration(registration)
	if err != nil {
		t.Fatalf("GetByRegistration failed: %v", err)
	}

	if data.Registration != registration {
		t.Errorf("Expected registration %s, got %s", registration, data.Registration)
	}
}

func TestGetByVIN(t *testing.T) {
	mockServer := createMockServer()
	defer mockServer.Close()

	BaseURL = mockServer.URL
	client := createTestClient(mockServer)

	registration := "ML58FOU"
	vin := "BNR32305366"
	data, err := client.GetByVIN(vin)
	if err != nil {
		t.Fatalf("GetByVIN failed: %v", err)
	}

	if data.Registration != registration {
		t.Errorf("Expected registration %s, got %s", registration, data.Registration)
	}
}

func TestGetBulkDownload(t *testing.T) {
	mockServer := createMockServer()
	defer mockServer.Close()

	BaseURL = mockServer.URL
	client := createTestClient(mockServer)

	data, err := client.GetBulkDownload()
	if err != nil {
		t.Fatalf("GetBulkDownload failed: %v", err)
	}

	if data.Bulk == nil {
		t.Error("Expected 'bulk' key in response")
	}
	if data.Bulk == nil {
		t.Error("Expected 'delta' key in response")
	}
}

func TestRenewCredentials(t *testing.T) {
	mockServer := createMockServer()
	defer mockServer.Close()

	BaseURL = mockServer.URL
	client := createTestClient(mockServer)

	email := "f@finbarrs.eu"
	apiKey := "dummy-api-key"
	data, err := client.RenewCredentials(apiKey, email)
	if err != nil {
		t.Fatalf("RenewCredentials failed: %v", err)
	}

	if len(data.ClientSecret) == 0 {
		t.Error("Expected 'clientSecret' length > 0")
	}
}

func TestInvalidCases(t *testing.T) {
	mockServer := createMockServer()
	defer mockServer.Close()

	BaseURL = mockServer.URL
	client := createTestClient(mockServer)

	t.Run("InvalidRegistration", func(t *testing.T) {
		_, err := client.GetByRegistration("INVALID")
		if err == nil {
			t.Fatal("Expected an error for invalid registration, but got none")
		}
	})

	t.Run("InvalidVIN", func(t *testing.T) {
		_, err := client.GetByVIN("INVALID")
		if err == nil {
			t.Fatal("Expected an error for invalid VIN, but got none")
		}
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		invalidClient := NewClient(ClientConfig{
			ClientID:     "invalid",
			ClientSecret: "invalid",
			APIKey:       "invalid",
		}, nil) // Use `nil` to use the MOT API endpoint

		_, err := invalidClient.GetByRegistration("ML58FOU")
		if err == nil {
			t.Fatal("Expected an error for invalid credentials, but got none")
		}
	})
}

func TestRateLimiting(t *testing.T) {
	mockServer := createMockServer()
	defer mockServer.Close()

	BaseURL = mockServer.URL
	client := createTestClient(mockServer)

	// Adjust limiter parameters for testing
	// !! Increasing these variables will increase time taken to complete test
	rps := 1
	burst := 4
	client.rateLimiter = *rate.NewLimiter(rate.Limit(rps), burst)

	for i := 0; i < (burst + 1); i++ {
		registration := "ML58FOU"
		_, err := client.GetByRegistration(registration)
		if err != nil {
			t.Fatalf("Error occurred when testing rate limiting: %v", err)
		}

		if i == burst && client.rateLimiter.Tokens() >= 1 {
			t.Fatal("Rate limiting failed. After Burst tokens expected < 1")
		}
	}

	time.Sleep(1 * time.Second)
	if client.rateLimiter.Tokens() < 1 {
		t.Fatal("Rate limiting failed. After 1 second tokens expected > 1")
	}
}

func TestDayLimiting(t *testing.T) {
	mockServer := createMockServer()
	defer mockServer.Close()

	BaseURL = mockServer.URL
	client := createTestClient(mockServer)

	// Adjust limiter parameters for testing
	// !! Increasing these variables will increase time taken to complete test
	dailyQuota := 4
	secondsInDay := 10

	dailyRate := rate.Limit(float64(dailyQuota) / float64(secondsInDay))
	client.dayLimiter = *rate.NewLimiter(dailyRate, dailyQuota)

	for i := 0; i < (dailyQuota + 1); i++ {
		registration := "ML58FOU"
		_, err := client.GetByRegistration(registration)
		if err != nil {
			t.Fatalf("Error occurred when testing rate limiting: %v", err)
		}

		if i == dailyQuota && client.dayLimiter.Tokens() > 1 {
			t.Fatal("Day limiting failed. After using daily quota tokens expected < 1")
		}
	}
}

func createTestClient(mockServer *httptest.Server) *Client {
	mockConfig := ClientConfig{
		ClientID:     "nil",
		ClientSecret: "nil",
		APIKey:       "nil",
	}
	return NewClient(mockConfig, mockServer.Client())
}
