package client

import (
	"net/http"
	"time"
)

// Client represents an API Client for the triplestore API.
//
// It holds both functions specific to GraphDB, as well as generic HTTP functions.
type Client struct {
	// Client is a http.Client to connect to the api.
	Client http.Client

	// BaseURL is the base URL of the triplestore API.
	BaseURL string

	// Username and Password hold credentials for authenticating.
	Username string
	Password string

	// Duration used by wait to wait for operations
	PollInterval time.Duration
}

func NewClient(timeout time.Duration, baseURL, adminUsername, adminPassword string) *Client {
	return &Client{
		Client: http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				DisableKeepAlives: true,
			},
		},
		BaseURL:  baseURL,
		Username: adminUsername,
		Password: adminPassword,

		PollInterval: 1 * time.Second,
	}
}

// requestHeaders represent headers of a raw http request.
type requestHeaders struct {
	Accept      string
	ContentType string
}

func (rh *requestHeaders) With(headers requestHeaders) *requestHeaders {
	// create new request headers and copy the old options
	var newHeaders requestHeaders
	if rh != nil {
		newHeaders = *rh
	}

	// add the options
	if headers.Accept != "" {
		newHeaders.Accept = headers.Accept
	}

	if headers.ContentType != "" {
		newHeaders.ContentType = headers.ContentType
	}

	return &newHeaders
}
