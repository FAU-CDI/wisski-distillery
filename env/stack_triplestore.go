package env

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
	"github.com/FAU-CDI/wisski-distillery/internal/wait"
	"github.com/pkg/errors"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

func (dis *Distillery) TriplestoreStack() stack.Installable {
	return dis.asCoreStack("triplestore", stack.Installable{
		CopyContextFiles: []string{"graphdb.zip"},

		MakeDirsPerm: fs.ModeDir | fs.ModePerm,
		MakeDirs: []string{
			filepath.Join("data", "data"),
			filepath.Join("data", "work"),
			filepath.Join("data", "logs"),
		},
	})
}

func (dis *Distillery) TriplestoreStackPath() string {
	return dis.TriplestoreStack().Dir
}

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

const triplestoreBaseURL = "http://127.0.0.1:7200"
const waitTSInterval = 1 * time.Second

// triplestoreCall makes a request to the triplestore.
//
// When bodyName is non-empty, expect body to be a byte slice representing a multipart/form-data upload with the given name.
// When bodyName is empty, simply marshal body as application/json
func (dis *Distillery) triplestoreRequest(method, url string, body interface{}, bodyName string, accept string) (*http.Response, error) {
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
	req, err := http.NewRequest(method, triplestoreBaseURL+url, reader)
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
	req.SetBasicAuth(dis.Config.TriplestoreAdminUser, dis.Config.TriplestoreAdminPassword)

	// and send it
	return http.DefaultClient.Do(req)
}

func (dis *Distillery) TriplestoreWaitForConnection() error {
	return wait.Wait(func() bool {
		res, err := dis.triplestoreRequest("GET", "/rest/repositories", nil, "", "")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return true
	}, waitTSInterval, dis.Context())
}

var errTripleStoreFailedRepository = exit.Error{
	Message:  "Failed to create repository: %s",
	ExitCode: exit.ExitGeneric,
}

func (dis *Distillery) TriplestoreProvision(name, domain, user, password string) error {
	if err := dis.TriplestoreWaitForConnection(); err != nil {
		return err
	}

	// prepare the create repo request
	createRepo, err := distillery.ReadTemplate(filepath.Join("resources", "templates", "repository", "graphdb-repo.ttl"), map[string]string{
		"GRAPHDB_REPO":    name,
		"INSTANCE_DOMAIN": domain,
	})
	if err != nil {
		return err
	}

	// do the create!
	{
		res, err := dis.triplestoreRequest("POST", "/rest/repositories", createRepo, "config", "")
		if err != nil {
			return errTripleStoreFailedRepository.WithMessageF(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusCreated {
			return errTripleStoreFailedRepository.WithMessageF("Repo create did not return status code 201")
		}
	}

	// create the user and grant them access
	{
		res, err := dis.triplestoreRequest("POST", "/rest/security/users/"+user, TriplestoreUserPayload{
			Password: password,
			AppSettings: TriplestoreUserAppSettings{
				DefaultInference:      true,
				DefaultVisGraphSchema: true,
				DefaultSameas:         true,
				IgnoreSharedQueries:   false,
				ExecuteCount:          true,
			},
			GrantedAuthorities: []string{
				"ROLE_USER",
				"READ_REPO_" + name,
				"WRITE_REPO_" + name,
			},
		}, "", "")
		if err != nil {
			return errTripleStoreFailedRepository.WithMessageF(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusCreated {
			return errTripleStoreFailedRepository.WithMessageF("User create did not return status code 201")
		}
	}

	return nil
}

// TriplestorePurgeUser deletes the specified user from the triplestore
func (dis *Distillery) TriplestorePurgeUser(user string) error {
	res, err := dis.triplestoreRequest("DELETE", "/rest/security/users/"+user, nil, "", "")
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusNoContent {
		return errors.Errorf("Delete returned code %d", res.StatusCode)
	}
	return nil
}

// TriplestorePurgeRepo deletes the specified repo from the triplestore
func (dis *Distillery) TriplestorePurgeRepo(repo string) error {
	res, err := dis.triplestoreRequest("DELETE", "/rest/repositories/"+repo, nil, "", "")
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.Errorf("Delete returned code %d", res.StatusCode)
	}
	return nil
}

var errTriplestoreFailedSecurity = errors.New("failed to enable triplestore security: request did not succeed with HTTP 200 OK")

func (dis *Distillery) TriplestoreBootstrap(io stream.IOStream) error {
	logging.LogMessage(io, "Waiting for Triplestore")
	if err := dis.TriplestoreWaitForConnection(); err != nil {
		return err
	}

	logging.LogMessage(io, "Resetting admin user password")
	{
		res, err := dis.triplestoreRequest("PUT", "/rest/security/users/"+dis.Config.TriplestoreAdminUser, TriplestoreUserPayload{
			Password: dis.Config.TriplestoreAdminPassword,
			AppSettings: TriplestoreUserAppSettings{
				DefaultInference:      true,
				DefaultVisGraphSchema: true,
				DefaultSameas:         true,
				IgnoreSharedQueries:   false,
				ExecuteCount:          true,
			},
			GrantedAuthorities: []string{"ROLE_ADMIN"},
		}, "", "")
		if err != nil {
			return fmt.Errorf("failed to create triplestore user: %s", err)
		}
		defer res.Body.Close()

		switch res.StatusCode {
		case http.StatusOK:
			// we set the password => requests are unauthorized
			// so we still need to enable security (see below!)
		case http.StatusUnauthorized:
			// a password is needed => security is already enabled.
			// the password may or may not work, but that's a problem for later
			logging.LogMessage(io, "Security is already enabled")
			return nil
		default:
			return fmt.Errorf("failed to create triplestore user: %s", err)
		}
	}

	logging.LogMessage(io, "Enabling Triplestore security")
	{
		res, err := dis.triplestoreRequest("POST", "/rest/security", true, "", "")
		if err != nil {
			return fmt.Errorf("failed to enable triplestore security: %s", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return errTriplestoreFailedSecurity
		}

		return nil
	}
}
