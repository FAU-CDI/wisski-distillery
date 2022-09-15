package instances

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/stringparser"
	"github.com/alessio/shellescape"
	"github.com/tkw1536/goprogram/stream"
)

var errInvalidSlug = errors.New("not a valid slug")

// Create fills the struct for a new WissKI instance.
// It validates that slug is a valid name for an instance.
//
// It does not perform any checks if the instance already exists, or does the creation in the database.
func (instances *Instances) Create(slug string) (wisski WissKI, err error) {

	// make sure that the slug is valid!
	if _, err := stringparser.ParseSlug(slug); err != nil {
		return wisski, errInvalidSlug
	}

	wisski.Instance.Slug = slug
	wisski.Instance.FilesystemBase = filepath.Join(instances.Dir, slug)

	wisski.Instance.OwnerEmail = ""
	wisski.Instance.AutoBlindUpdateEnabled = true

	// sql

	wisski.Instance.SqlDatabase = instances.Config.MysqlDatabasePrefix + slug
	wisski.Instance.SqlUsername = instances.Config.MysqlUserPrefix + slug

	wisski.Instance.SqlPassword, err = instances.Config.NewPassword()
	if err != nil {
		return WissKI{}, err
	}

	// triplestore

	wisski.Instance.GraphDBRepository = instances.Config.GraphDBRepoPrefix + slug
	wisski.Instance.GraphDBUsername = instances.Config.GraphDBUserPrefix + slug

	wisski.Instance.GraphDBPassword, err = instances.Config.NewPassword()
	if err != nil {
		return WissKI{}, err
	}

	// drupal

	wisski.DrupalUsername = "admin" // TODO: Change this!

	wisski.DrupalPassword, err = instances.Config.NewPassword()
	if err != nil {
		return wisski, err
	}

	// store the instance in the object and return it!
	wisski.instances = instances
	return wisski, nil
}

// Provision provisions an instance, assuming that the required databases already exist.
func (wisski WissKI) Provision(io stream.IOStream) error {

	// create the basic st!
	st := wisski.Barrel()
	if err := st.Install(io, component.InstallationContext{}); err != nil {
		return err
	}

	// Pull and build the stack!
	if err := st.Update(io, false); err != nil {
		return err
	}

	provisionParams := []string{
		wisski.Domain(),

		wisski.SqlDatabase,
		wisski.SqlUsername,
		wisski.SqlPassword,

		wisski.GraphDBRepository,
		wisski.GraphDBUsername,
		wisski.GraphDBPassword,

		wisski.DrupalUsername,
		wisski.DrupalPassword,

		"", // TODO: DrupalVersion
		"", // TODO: WissKIVersion
	}

	// escape the parameter
	for i, param := range provisionParams {
		provisionParams[i] = shellescape.Quote(param)
	}

	// figure out the provision script
	// TODO: Move the provision script into the control plane!
	provisionScript := "sudo PATH=$PATH -u www-data /bin/bash /provision_container.sh " + strings.Join(provisionParams, " ")

	code, err := st.Run(io, true, "barrel", "/bin/bash", "-c", provisionScript)
	if err != nil {
		return err
	}
	if code != 0 {
		return errors.New("unable to run provision script")
	}

	return nil
}
