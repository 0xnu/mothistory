package mothistory

import (
	"encoding/json"
	"os"
	"testing"
)

func TestMOTHistoryClient(t *testing.T) {
	clientID := os.Getenv("MOT_CLIENT_ID")
	clientSecret := os.Getenv("MOT_CLIENT_SECRET")
	apiKey := os.Getenv("MOT_API_KEY")

	if clientID == "" || clientSecret == "" || apiKey == "" {
		t.Fatal("Environment variables MOT_CLIENT_ID, MOT_CLIENT_SECRET, and MOT_API_KEY must be set")
	}

	config := ClientConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		APIKey:       apiKey,
	}

	client := NewClient(config)

	t.Run("GetByRegistration", func(t *testing.T) {
		registration := "ML58FOU" // Please replace with a valid registration number
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
		vin := "BNR32305366" // Please replace with a valid VIN
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

	t.Run("InvalidCredentials", func(t *testing.T) {
		invalidClient := NewClient(ClientConfig{
			ClientID:     "invalid",
			ClientSecret: "invalid",
			APIKey:       "invalid",
		})

		_, err := invalidClient.GetByRegistration("ML58FOU")
		if err == nil {
			t.Fatal("Expected an error for invalid credentials, but got none")
		}
	})
}
