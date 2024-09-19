package mothistory

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	// "os"
	"testing"
)

func createMockAPI() (*httptest.Server) {
	handler := http.NewServeMux()

	handler.HandleFunc("/registration/ML58FOU", func(w http.ResponseWriter, r *http.Request) {
		mockResponse := `{"registration": "ML58FOU"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockResponse)
	})

	handler.HandleFunc("/vin/BNR32305366", func(w http.ResponseWriter, r *http.Request) {
		mockResponse := `{"vin": "BNR32305366"}`
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

func TestMOTHistoryClient(t *testing.T) {
	mockServer := createMockAPI()
	defer mockServer.Close()

	BaseURL = mockServer.URL

	mockConfig := ClientConfig{
		ClientID:     "nil",
		ClientSecret: "nil",
		APIKey:       "nil",
	}

	client := NewClient(mockConfig, mockServer.Client())

	t.Run("GetByRegistration", func(t *testing.T) {
		registration := "ML58FOU"
		data, err := client.GetByRegistration(registration)
		if err != nil {
			t.Fatalf("GetByRegistration failed: %v", err)
		}

		var response map[string]interface{}
		err = json.Unmarshal(data, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response["registration"] != registration {
			t.Errorf("Expected registration %s, got %s", registration, response["registration"])
		}
	})

	t.Run("GetByVIN", func(t *testing.T) {
		vin := "BNR32305366"
		data, err := client.GetByVIN(vin)
		if err != nil {
			t.Fatalf("GetByVIN failed: %v", err)
		}

		var response map[string]interface{}
		err = json.Unmarshal(data, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response["vin"] != vin {
			t.Errorf("Expected VIN %s, got %s", vin, response["vin"])
		}
	})

	t.Run("GetBulkDownload", func(t *testing.T) {
		data, err := client.GetBulkDownload()
		if err != nil {
			t.Fatalf("GetBulkDownload failed: %v", err)
		}

		var response map[string]interface{}
		err = json.Unmarshal(data, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if _, ok := response["bulk"]; !ok {
			t.Error("Expected 'bulk' key in response")
		}
		if _, ok := response["delta"]; !ok {
			t.Error("Expected 'delta' key in response")
		}
	})

	t.Run("RenewCredentials", func(t *testing.T) {
		email := "f@finbarrs.eu" // Please replace your email
		apiKey := "dummy-api-key"
		data, err := client.RenewCredentials(apiKey, email)
		if err != nil {
			t.Fatalf("RenewCredentials failed: %v", err)
		}

		var response map[string]interface{}
		err = json.Unmarshal(data, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if _, ok := response["clientSecret"]; !ok {
			t.Error("Expected 'clientSecret' key in response")
		}
	})

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

	// t.Run("InvalidCredentials", func(t *testing.T) {
	// 	invalidClient := NewClient(ClientConfig{
	// 		ClientID:     "invalid",
	// 		ClientSecret: "invalid",
	// 		APIKey:       "invalid",
	// 	})

	// 	_, err := invalidClient.GetByRegistration("ML58FOU")
	// 	if err == nil {
	// 		t.Fatal("Expected an error for invalid credentials, but got none")
	// 	}
	// })
}
