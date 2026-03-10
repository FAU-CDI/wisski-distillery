package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
	"strings"
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

type headers struct {
	Accept      string
	ContentType string
}

// rest performs a (raw) http request to the triplestore API.
func (client *Client) rest(ctx context.Context, method, url string, h headers) (*http.Response, error) {
	return client.doRestWithReader(ctx, method, url, h, nil)
}

// doRestWithForm performs a http request where the body are all bytes read from fieldvalue.
func (client *Client) doRestWithForm(ctx context.Context, method, url string, h headers, fieldname string, fieldvalue io.Reader) (*http.Response, error) {
	var buffer bytes.Buffer

	// write the file to it
	writer := multipart.NewWriter(&buffer)
	{
		part, err := writer.CreateFormFile(fieldname, "filename.txt")
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		if _, err := io.Copy(part, fieldvalue); err != nil {
			return nil, fmt.Errorf("failed to copy values into form: %w", err)
		}
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// set it with the given content type
	h.ContentType = writer.FormDataContentType()
	return client.doRestWithReader(ctx, method, url, h, &buffer)
}

// DoRestWithReader performs a http request where the body is copied from the given io.Reader.
// The caller must ensure the reader is closed.
func (client *Client) doRestWithMarshal(ctx context.Context, method, url string, h headers, body any) (*http.Response, error) {
	// encode into a buffer
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(body); err != nil {
		return nil, fmt.Errorf("failed to encode body: %w", err)
	}

	h.ContentType = "application/json"
	return client.doRestWithReader(ctx, method, url, h, &buffer)
}

// doRestWithReader performs a http request where the body is copied from the given io.Reader.
// The caller must ensure the reader is closed.
func (ts *Client) doRestWithReader(ctx context.Context, method string, url string, h headers, body io.Reader) (*http.Response, error) {
	// create the request and authentication
	req, err := http.NewRequestWithContext(ctx, method, ts.BaseURL+url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(ts.Username, ts.Password)

	// add extra headers
	if h.Accept != "" {
		req.Header.Set("Accept", h.Accept)
	}
	if h.ContentType != "" {
		req.Header.Set("Content-Type", h.ContentType)
	}

	// and send it
	res, err := ts.Client.Do(req)
	if err != nil {
		return res, fmt.Errorf("(http.Client).Do reported: %w", err)
	}
	return res, nil
}

type WrongStatusError struct {
	Expected []int
	Got      int
	Body     string
}

func (e WrongStatusError) Error() string {
	expected := make([]string, len(e.Expected))
	for i, code := range e.Expected {
		expected[i] = strconv.Itoa(code)
	}
	expectedString := strings.Join(expected, ", ")
	if len(expected) > 1 {
		expectedString = "one of " + expectedString
	}

	if e.Body != "" {
		return fmt.Sprintf("wrong status code: expected %s, got %d: %s", expectedString, e.Got, e.Body)
	}
	return fmt.Sprintf("wrong status code: expected %s, got %d", expectedString, e.Got)
}

// newStatusError checks that the response code of res matches, and returns an [WrongStatusError] (or an error that contains it).
// If withBody is true, the body of the response is read and included in the error message.
func newStatusError(res *http.Response, codes ...int) error {
	if res == nil || slices.Contains(codes, res.StatusCode) {
		return nil
	}

	body, bodyErr := io.ReadAll(res.Body)

	result := WrongStatusError{Expected: codes, Got: res.StatusCode, Body: string(body)}
	if bodyErr != nil {
		return errors.Join(
			result,
			fmt.Errorf("failed to read body for error message: %w", bodyErr),
		)
	}
	return result
}
