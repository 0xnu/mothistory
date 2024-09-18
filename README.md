## MOT History Go SDK

[![Release](https://img.shields.io/github/release/0xnu/mothistory.svg)](https://github.com/0xnu/mothistory/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xnu/mothistory)](https://goreportcard.com/report/github.com/0xnu/mothistory)
[![Go Reference](https://pkg.go.dev/badge/github.com/0xnu/mothistory.svg)](https://pkg.go.dev/github.com/0xnu/mothistory)
[![License](https://img.shields.io/github/license/0xnu/mothistory)](/LICENSE)

The SDK provides convenient access to the [MOT History API](https://documentation.history.mot.api.gov.uk/) for applications written in the [Go](https://go.dev/) Programming Language.

### Tests

Export environment variables:

```sh
export MOT_CLIENT_ID=
export MOT_CLIENT_SECRET=
export MOT_API_KEY=
```

Now, you can execute this command: `go test -v`

### Integration Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	mothistory "github.com/0xnu/mothistory"
)

const (
	clientID     = "enter_your_client_id>"
	clientSecret = "enter_your_client_secret>"
	apiKey       = "enter_your_api_key"
)

func main() {
	config := mothistory.ClientConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		APIKey:       apiKey,
	}
	client := mothistory.NewClient(config)

	// Get data by registration
	data, err := client.GetByRegistration("ML58FOU")
	if err != nil {
		log.Fatalf("failed to get data by registration: %v", err)
	}
	printJSON(data)

	// Get data by VIN
	data, err = client.GetByVIN("AISXXXTEST1239617")
	if err != nil {
		log.Fatalf("failed to get data by VIN: %v", err)
	}
	printJSON(data)

	// Get bulk download data
	data, err = client.GetBulkDownload()
	if err != nil {
		log.Fatalf("failed to get bulk download data: %v", err)
	}
	printJSON(data)

	// Renew credentials
	data, err = client.RenewCredentials(apiKey, "firstname.lastname@example.com")
	if err != nil {
		log.Fatalf("failed to renew credentials: %v", err)
	}
	printJSON(data)
}

func printJSON(data json.RawMessage) {
	var prettyJSON map[string]interface{}
	err := json.Unmarshal(data, &prettyJSON)
	if err != nil {
		log.Fatalf("failed to parse JSON: %v", err)
	}

	prettyData, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		log.Fatalf("failed to format JSON: %v", err)
	}

	fmt.Println(string(prettyData))
}
```

### Setting up a MOT History API

You can use this support form to request an [API Key](https://documentation.history.mot.api.gov.uk/mot-history-api/register).


### Using the MOT History API

You can read the [API documentation](https://documentation.history.mot.api.gov.uk/) to understand what's possible with the MOT History API. If you need further assistance, don't hesitate to [contact the DVSA](https://documentation.history.mot.api.gov.uk/mot-history-api/support).


### License

This project is licensed under the [MIT License](./LICENSE).


### Copyright

(c) 2023 - 2024 [Finbarrs Oketunji](https://finbarrs.eu).

The MOT History API Go SDK is Licensed under the [Open Government Licence v3.0](
https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/)
