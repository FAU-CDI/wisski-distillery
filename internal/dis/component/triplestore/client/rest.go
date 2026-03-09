package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// rest performs a (raw) http request to the triplestore API.
func (client *Client) rest(ctx context.Context, method, url string, headers *requestHeaders) (*http.Response, error) {
	return client.doRestWithReader(ctx, method, url, headers, nil)
}

// doRestWithForm performs a http request where the body are all bytes read from fieldvalue.
func (client *Client) doRestWithForm(ctx context.Context, method, url string, headers *requestHeaders, fieldname string, fieldvalue io.Reader) (*http.Response, error) {
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

	// and sent the reader as the body
	return client.doRestWithReader(ctx, method, url, headers.With(requestHeaders{ContentType: writer.FormDataContentType()}), &buffer)
}

// DoRestWithReader performs a http request where the body is copied from the given io.Reader.
// The caller must ensure the reader is closed.
func (client *Client) doRestWithMarshal(ctx context.Context, method, url string, headers *requestHeaders, body any) (*http.Response, error) {
	// encode into a buffer
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(body); err != nil {
		return nil, fmt.Errorf("failed to encode body: %w", err)
	}

	return client.doRestWithReader(ctx, method, url, headers.With(requestHeaders{ContentType: "application/json"}), &buffer)
}

// doRestWithReader performs a http request where the body is copied from the given io.Reader.
// The caller must ensure the reader is closed.
func (ts *Client) doRestWithReader(ctx context.Context, method string, url string, headers *requestHeaders, body io.Reader) (*http.Response, error) {

	// create the request and authentication
	req, err := http.NewRequestWithContext(ctx, method, ts.BaseURL+url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request: %w", err)
	}
	req.SetBasicAuth(ts.Username, ts.Password)

	// add extra headers
	if headers != nil && headers.Accept != "" {
		req.Header.Set("Accept", headers.Accept)
	}
	if headers != nil && headers.ContentType != "" {
		req.Header.Set("Content-Type", headers.ContentType)
	}

	// and send it
	res, err := ts.Client.Do(req)
	if err != nil {
		return res, fmt.Errorf("failed to do http request: %w", err)
	}
	return res, nil
}
