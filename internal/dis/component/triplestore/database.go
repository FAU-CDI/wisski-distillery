//spellchecker:words triplestore
package triplestore

//spellchecker:words bytes context encoding json mime multipart http time github wisski distillery internal component wdlog errors pkglib timex
import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/pkg/errors"
	"github.com/tkw1536/pkglib/timex"
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

// http.Client Timeout to be used for "trivial" triplestore operations.
// This includes e.g. CRUDing a specific repo.
const tsTrivialTimeout = time.Minute

// RequestHeaders represent headers of a raw http request
type RequestHeaders struct {
	Accept      string
	ContentType string
}

func (rh *RequestHeaders) With(headers RequestHeaders) *RequestHeaders {

	// create new request headers and copy the old options
	var newHeaders RequestHeaders
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

// DoRest performs a (raw) http request to the without a body.
func (ts *Triplestore) DoRest(ctx context.Context, timeout time.Duration, method, url string, headers *RequestHeaders) (*http.Response, error) {
	return ts.DoRestWithReader(ctx, timeout, method, url, headers, nil)
}

// DoRestWithForm performs a http request where the body are all bytes read from fieldvalue.
func (ts *Triplestore) DoRestWithForm(ctx context.Context, timeout time.Duration, method, url string, headers *RequestHeaders, fieldname string, fieldvalue io.Reader) (*http.Response, error) {
	var buffer bytes.Buffer

	// write the file to it
	writer := multipart.NewWriter(&buffer)
	{
		part, err := writer.CreateFormFile(fieldname, "filename.txt")
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(part, fieldvalue); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// and sent the reader as the body
	return ts.DoRestWithReader(ctx, timeout, method, url, headers.With(RequestHeaders{ContentType: writer.FormDataContentType()}), &buffer)
}

// DoRestWithReader performs a http request where the body is copied from the given io.Reader.
// The caller must ensure the reader is closed.
func (ts *Triplestore) DoRestWithMarshal(ctx context.Context, timeout time.Duration, method, url string, headers *RequestHeaders, body any) (*http.Response, error) {
	// encode into a buffer
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(body); err != nil {
		return nil, err
	}

	return ts.DoRestWithReader(ctx, timeout, method, url, headers.With(RequestHeaders{ContentType: "application/json"}), &buffer)
}

// DoRestWithReader performs a http request where the body is copied from the given io.Reader.
// The caller must ensure the reader is closed.
func (ts *Triplestore) DoRestWithReader(ctx context.Context, timeout time.Duration, method string, url string, headers *RequestHeaders, body io.Reader) (*http.Response, error) {
	// create the request object
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	config := component.GetStill(ts).Config.TS

	// create the request and authentication
	req, err := http.NewRequestWithContext(ctx, method, ts.BaseURL+url, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(config.AdminUsername, config.AdminPassword)

	// add extra headers
	if headers != nil && headers.Accept != "" {
		req.Header.Set("Accept", headers.Accept)
	}
	if headers != nil && headers.ContentType != "" {
		req.Header.Set("Content-Type", headers.ContentType)
	}

	// and send it
	return client.Do(req)
}

// Wait waits for the connection to the Triplestore to succeed.
// This is achieved using a polling strategy.
func (ts Triplestore) Wait(ctx context.Context) error {
	return timex.TickUntilFunc(func(time.Time) bool {
		res, err := ts.DoRest(ctx, tsTrivialTimeout, http.MethodGet, "/rest/repositories", nil)
		wdlog.Of(ctx).Debug(
			"Triplestore Wait",
			"error", err,
		)
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
	res, err := ts.DoRest(ctx, tsTrivialTimeout, http.MethodDelete, "/rest/security/users/"+user, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusNotFound {
		return errors.Errorf("Delete returned code %d", res.StatusCode)
	}
	return nil
}

// PurgeRepo deletes the specified repo from the triplestore.
// When the repo does not exist, returns no error.
func (ts Triplestore) PurgeRepo(ctx context.Context, repo string) error {
	res, err := ts.DoRest(ctx, tsTrivialTimeout, http.MethodDelete, "/rest/repositories/"+repo, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
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
