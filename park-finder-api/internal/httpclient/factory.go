package httpclient

import (
	"fmt"
	"net/http"
)

type HTTPClient struct {
	Client http.Client
}

func NewHTTPClient(c http.Client) *HTTPClient {
	return &HTTPClient{
		Client: c,
	}
}

// formatBearerToken format a string in to Bearer 'token' format
func formatBearerToken(token string) string {
	return fmt.Sprintf("Bearer %v", token)
}
