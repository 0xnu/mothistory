## MOT History Go SDK

[![Release](https://img.shields.io/github/release/0xnu/mothistory.svg)](https://github.com/0xnu/mothistory/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/0xnu/mothistory)](https://goreportcard.com/report/github.com/0xnu/mothistory)
[![Go Reference](https://pkg.go.dev/badge/github.com/0xnu/mothistory.svg)](https://pkg.go.dev/github.com/0xnu/mothistory)
[![License](https://img.shields.io/github/license/0xnu/mothistory)](/LICENSE)

The SDK provides convenient access to the [MOT History API](https://documentation.history.mot.api.gov.uk/) for applications written in the [Go](https://go.dev/) Programming Language.

### Usage Example

```go
package main

import (
	"fmt"

	mothistory "github.com/0xnu/mothistory"
)

const apiKey = "<your-api-key>"

func main() {
	client := mothistory.NewClient(apiKey)

	// Get data by registration
	json, err := client.GetByRegistration("ML58FOU")
	if err != nil {
		fmt.Printf("failed to get data by registration: %v", err)
		return
	}
	fmt.Println(json)

    // Get data by page
    json, err = client.GetByPage(2)
    if err != nil {
        fmt.Printf("failed to get data by page: %v", err)
        return
    }
    fmt.Println(json)

    // Get data by date and page
    json, err = client.GetByDateAndPage("20230201", 1)
    if err != nil {
        fmt.Printf("failed to get data by date and page: %v", err)
        return
    }
    fmt.Println(json)

    // Get data by vehicle ID
    json, err = client.GetByVehicleID("<enter your vehicle id here>")
    if err != nil {
        fmt.Printf("failed to get data by vehicle ID: %v", err)
        return
    }
    fmt.Println(json)

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
