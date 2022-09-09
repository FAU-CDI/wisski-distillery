package env

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/FAU-CDI/wisski-distillery/embed"
	"github.com/FAU-CDI/wisski-distillery/internal/logging"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
	"github.com/FAU-CDI/wisski-distillery/internal/wait"
	"github.com/pkg/errors"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

// TriplestoreComponent represents the triplestore belonging to a distillery
type TriplestoreComponent struct {
	BaseURL      string        // the base url of the api
	PollInterval time.Duration // duration to wait during wait!

	dis *Distillery
}

// Triplestore returns the TriplestoreComponent belonging to this distillery
func (dis *Distillery) Triplestore() TriplestoreComponent {
	return TriplestoreComponent{
		BaseURL:      "http://" + dis.Upstream.Triplestore,
		PollInterval: time.Second,

		dis: dis,
	}
}

func (TriplestoreComponent) Name() string {
	return "triplestore"
}

func (TriplestoreComponent) Context(parent stack.InstallationContext) stack.InstallationContext {
	return parent
}

// Stack returns the installable Triplestore stack
func (ts TriplestoreComponent) Stack() stack.Installable {
	return ts.dis.makeComponentStack(ts, stack.Installable{
		CopyContextFiles: []string{"graphdb.zip"},

		MakeDirsPerm: fs.ModeDir | fs.ModePerm,
		MakeDirs: []string{
			filepath.Join("data", "data"),
			filepath.Join("data", "work"),
			filepath.Join("data", "logs"),
		},
	})
}

func (ts TriplestoreComponent) Path() string {
	return ts.Stack().Dir
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

// OpenRaw makes an http request to the triplestore api.
//
// When bodyName is non-empty, expect body to be a byte slice representing a multipart/form-data upload with the given name.
// When bodyName is empty, simply marshal body as application/json
func (ts TriplestoreComponent) OpenRaw(method, url string, body interface{}, bodyName string, accept string) (*http.Response, error) {
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
	req, err := http.NewRequest(method, ts.BaseURL+url, reader)
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
	req.SetBasicAuth(ts.dis.Config.TriplestoreAdminUser, ts.dis.Config.TriplestoreAdminPassword)

	// and send it
	return http.DefaultClient.Do(req)
}

// Wait waits for the connection to the Triplestore to succeed.
// This is achieved using a polling strategy.
func (ts TriplestoreComponent) Wait() error {
	return wait.Wait(func() bool {
		res, err := ts.OpenRaw("GET", "/rest/repositories", nil, "", "")
		if err != nil {
			return false
		}
		defer res.Body.Close()
		return true
	}, ts.PollInterval, ts.dis.Context())
}

var errTripleStoreFailedRepository = exit.Error{
	Message:  "Failed to create repository: %s",
	ExitCode: exit.ExitGeneric,
}

func (ts TriplestoreComponent) Provision(name, domain, user, password string) error {
	if err := ts.Wait(); err != nil {
		return err
	}

	// prepare the create repo request
	createRepo, err := embed.ReadTemplate(filepath.Join("resources", "templates", "repository", "graphdb-repo.ttl"), map[string]string{
		"GRAPHDB_REPO":    name,
		"INSTANCE_DOMAIN": domain,
	})
	if err != nil {
		return err
	}

	// do the create!
	{
		res, err := ts.OpenRaw("POST", "/rest/repositories", createRepo, "config", "")
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
		res, err := ts.OpenRaw("POST", "/rest/security/users/"+user, TriplestoreUserPayload{
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
func (ts TriplestoreComponent) PurgeUser(user string) error {
	res, err := ts.OpenRaw("DELETE", "/rest/security/users/"+user, nil, "", "")
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusNoContent {
		return errors.Errorf("Delete returned code %d", res.StatusCode)
	}
	return nil
}

// TriplestorePurgeRepo deletes the specified repo from the triplestore
func (ts TriplestoreComponent) PurgeRepo(repo string) error {
	res, err := ts.OpenRaw("DELETE", "/rest/repositories/"+repo, nil, "", "")
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.Errorf("Delete returned code %d", res.StatusCode)
	}
	return nil
}

var errTSBackupWrongStatusCode = errors.New("Distillery.Backup: Wrong status code")

// TriplestoreBackup backs up the repository named repo into the writer dst.
func (ts TriplestoreComponent) Backup(dst io.Writer, repo string) (int64, error) {
	res, err := ts.OpenRaw("GET", "/repositories/"+repo+"/statements?infer=false", nil, "", "application/n-quads")
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, errTSBackupWrongStatusCode
	}
	defer res.Body.Close()
	return io.Copy(dst, res.Body)
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

func (ts TriplestoreComponent) listRepositories() (repos []Repository, err error) {
	res, err := ts.OpenRaw("GET", "/rest/repositories", nil, "", "application/json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&repos)
	return
}

// TriplestoreBackup backs up every graphdb instance into dst
func (ts TriplestoreComponent) BackupAll(dst string) error {
	// list all the repositories
	repos, err := ts.listRepositories()
	if err != nil {
		return err
	}

	// create the base directory
	if err := os.Mkdir(dst, fs.ModeDir); err != nil {
		return err
	}

	// iterate over all the repositories
	for _, repo := range repos {
		if rErr := (func(repo Repository) error {
			name := filepath.Join(dst, repo.ID+".nq")

			dest, err := os.Create(name)
			if err != nil {
				return err
			}
			defer dest.Close()

			_, err = ts.Backup(dest, repo.ID)
			return err
		}(repo)); err == nil && rErr != nil {
			err = rErr
		}
	}
	return err
}

var errTriplestoreFailedSecurity = errors.New("failed to enable triplestore security: request did not succeed with HTTP 200 OK")

func (ts TriplestoreComponent) Bootstrap(io stream.IOStream) error {
	logging.LogMessage(io, "Waiting for Triplestore")
	if err := ts.Wait(); err != nil {
		return err
	}

	logging.LogMessage(io, "Resetting admin user password")
	{
		res, err := ts.OpenRaw("PUT", "/rest/security/users/"+ts.dis.Config.TriplestoreAdminUser, TriplestoreUserPayload{
			Password: ts.dis.Config.TriplestoreAdminPassword,
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
		res, err := ts.OpenRaw("POST", "/rest/security", true, "", "")
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
