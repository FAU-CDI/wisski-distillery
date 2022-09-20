package triplestore

import (
	"bytes"
	"net/http"

	_ "embed"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/unpack"
	"github.com/tkw1536/goprogram/exit"
)

var errTripleStoreFailedRepository = exit.Error{
	Message:  "Failed to create repository: %s",
	ExitCode: exit.ExitGeneric,
}

//go:embed create-repo.ttl
var createRepoTTL []byte

func (ts *Triplestore) Provision(instance models.Instance, domain string) error {
	return ts.CreateRepository(instance.GraphDBRepository, domain, instance.GraphDBUsername, instance.GraphDBPassword)
}

func (ts *Triplestore) CreateRepository(name, domain, user, password string) error {
	if err := ts.Wait(); err != nil {
		return err
	}

	// prepare the create repo request
	var createRepo bytes.Buffer
	err := unpack.WriteTemplate(&createRepo, map[string]string{
		"GRAPHDB_REPO":    name,
		"INSTANCE_DOMAIN": domain,
	}, bytes.NewReader(createRepoTTL))
	if err != nil {
		return err
	}

	// do the create!
	{
		res, err := ts.OpenRaw("POST", "/rest/repositories", createRepo.Bytes(), "config", "")
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
