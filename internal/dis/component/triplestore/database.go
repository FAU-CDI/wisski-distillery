package triplestore

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/FAU-CDI/wisski-distillery/pkg/timex"
	"github.com/pkg/errors"
	"github.com/tkw1536/goprogram/stream"
)

type TriplestoreUserPayload struct {
	Password           string                     `json:"password"`
	AppSettings        TriplestoreUserAppSettings `json:"appSettings"`
	GrantedAuthorities []string                   `json:"grantedAuthorities"`
}
type TriplestoreUserAppSettings struct {
	DefaultInference      bool `json:"DEFAULT_INFERENCE"`
	DefaultVisGraphSchema bool `json:"DEFAULT_VIS_GRAPH_SCHEMA"`
	DefaultSameas         bool `json:"DEFAULT_SAMEAS"`
	IgnoreSharedQueries   bool `json:"IGNORE_SHARED_QUERIES"`
	ExecuteCount          bool `json:"EXECUTE_COUNT"`
}

// OpenRaw makes an http request to the triplestore api.
//
// When bodyName is non-empty, expect body to be a byte slice representing a multipart/form-data upload with the given name.
// When bodyName is empty, simply marshal body as application/json
func (ts Triplestore) OpenRaw(ctx context.Context, method, url string, body interface{}, bodyName string, accept string) (*http.Response, error) {
	var reader io.Reader

	var contentType string

	// for "PUT" and "POST" we setup a body
	if method == "PUT" || method == "POST" {
		if bodyName != "" {
			buffer := &bytes.Buffer{}
			writer := multipart.NewWriter(buffer)
			contentType = writer.FormDataContentType()

			part, err := writer.CreateFormFile(bodyName, "filename.txt")
			if err != nil {
				return nil, err
			}
			io.Copy(part, bytes.NewReader(body.([]byte)))
			writer.Close()
			reader = buffer
		} else {
			contentType = "application/json"
			mbytes, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reader = bytes.NewReader(mbytes)
		}
	}

	// create the request object
	client := &http.Client{
		Transport: &http.Transport{
			DialContext:       ts.Environment.DialContext,
			DisableKeepAlives: true,
		},
	}
	req, err := http.NewRequestWithContext(ctx, method, ts.BaseURL+url, reader)
	if err != nil {
		return nil, err
	}

	// Setup configuration!
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.SetBasicAuth(ts.Config.TriplestoreAdminUser, ts.Config.TriplestoreAdminPassword)

	// and send it
	return client.Do(req)
}

// Wait waits for the connection to the Triplestore to succeed.
// This is achieved using a polling strategy.
func (ts Triplestore) Wait(ctx context.Context) error {
	n := stream.FromNil()
	return timex.TickUntilFunc(func(time.Time) bool {
		res, err := ts.OpenRaw(ctx, "GET", "/rest/repositories", nil, "", "")
		n.EPrintf("[Triplestore.Wait]: %s\n", err)
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return true
	}, ctx, ts.PollInterval)
}

// PurgeUser deletes the specified user from the triplestore.
// When the user does not exist, returns no error.
func (ts Triplestore) PurgeUser(ctx context.Context, user string) error {
	res, err := ts.OpenRaw(ctx, "DELETE", "/rest/security/users/"+user, nil, "", "")
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusNotFound {
		return errors.Errorf("Delete returned code %d", res.StatusCode)
	}
	return nil
}

// PurgeRepo deletes the specified repo from the triplestore.
// When the repo does not exist, returns no error.
func (ts Triplestore) PurgeRepo(ctx context.Context, repo string) error {
	res, err := ts.OpenRaw(ctx, "DELETE", "/rest/repositories/"+repo, nil, "", "")
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		return errors.Errorf("Delete returned code %d", res.StatusCode)
	}
	return nil
}

type Repository struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	URI        string `json:"uri"`
	Type       string `json:"type"`
	SesameType string `json:"sesameType"`
	Location   string `json:"location"`
	Readable   bool   `json:"readable"`
	Writable   bool   `json:"writable"`
	Local      bool   `json:"local"`
}
