package caller

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sfit-platform-web-backend/dtos"
)

// GetRequest makes a GET request to the specified URL with optional query parameters
func GetRequest(urlApi string, params map[string]string) (dtos.ApiCallerRp, error) {
	apiCallerRp := dtos.ApiCallerRp{}

	// Parse the URL
	apiURL, err := url.Parse(urlApi)
	if err != nil {
		apiCallerRp.StatusCode = 400
		return apiCallerRp, fmt.Errorf("error parsing URL: %w", err)
	}

	// Add query parameters if provided
	if params != nil {
		query := apiURL.Query()
		for key, value := range params {
			query.Add(key, value)
		}
		apiURL.RawQuery = query.Encode()
	}

	// Make the GET request
	rp, err := http.Get(apiURL.String())
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer rp.Body.Close()

	apiCallerRp.StatusCode = rp.StatusCode
	apiCallerRp.Body, _ = io.ReadAll(rp.Body)

	if rp.StatusCode != http.StatusOK {
		return apiCallerRp, fmt.Errorf("received status code %d", rp.StatusCode)
	}

	return apiCallerRp, nil
}
