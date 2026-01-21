package caller

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sfit-platform-web-backend/internal/dtos"
)

// GetRequest makes a GET request to the specified URL with optional query parameters
func GetRequest(urlApi string, params map[string]string) (dtos.ApiCallerRp, error) {
	apiCallerRp := dtos.ApiCallerRp{}

	// Parse URL
	apiURL, err := url.Parse(urlApi)
	if err != nil {
		return apiCallerRp, fmt.Errorf("error parsing URL: %w", err)
	}

	// Add query parameters
	if params != nil {
		query := apiURL.Query()
		for key, value := range params {
			query.Add(key, value)
		}
		apiURL.RawQuery = query.Encode()
	}

	// Make request
	rp, err := http.Get(apiURL.String())
	if err != nil {
		return apiCallerRp, fmt.Errorf("error making request: %w", err)
	}
	defer rp.Body.Close()

	apiCallerRp.StatusCode = rp.StatusCode

	// Read body
	body, err := io.ReadAll(rp.Body)
	if err != nil {
		return apiCallerRp, fmt.Errorf("read response body failed: %w", err)
	}
	apiCallerRp.Body = body

	// Check status code
	if rp.StatusCode != http.StatusOK {
		return apiCallerRp, fmt.Errorf(
			"received status code %d: %s",
			rp.StatusCode,
			string(body),
		)
	}

	return apiCallerRp, nil
}
