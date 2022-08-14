package env

import (
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/internal/password"
	"github.com/pkg/errors"
)

func (dis *Distillery) InstancesDir() string {
	return filepath.Join(dis.Config.DeployRoot, "instances")
}

func (dis *Distillery) InstanceDir(slug string) string {
	return filepath.Join(dis.InstancesDir(), slug)
}

func (dis *Distillery) InstanceSQL(slug string) (database, user string) {
	database = dis.Config.MysqlDatabasePrefix + slug
	user = dis.Config.MysqlUserPrefix + slug
	return
}

func (dis *Distillery) InstanceGraphDB(slug string) (repo, user string) {
	repo = dis.Config.GraphDBRepoPrefix + slug
	user = dis.Config.GraphDBUserPrefix + slug
	return
}

// Password returns a new password
func (dis *Distillery) NewPassword() (value string, err error) {
	return password.Password(dis.Config.PasswordLength)
}

var errInvalidSlug = errors.New("Not a valid slug")

// NewInstance fills the struct for a new distillery instance.
// It validates that slug is a valid name for an instance.
//
// It does not perform any checks if the instance already exists, or does the creation in the database.
func (dis *Distillery) NewInstance(slug string) (i Instance, err error) {

	// make sure that the slug is valid!
	if _, err := config.IsValidSlug(slug); err != nil {
		return i, errInvalidSlug
	}

	// generate sql data
	sqlPassword, err := dis.NewPassword()
	if err != nil {
		return i, err
	}
	sqlDB, sqlUser := dis.InstanceSQL(slug)

	// generate ts data
	tsPassword, err := dis.NewPassword()
	if err != nil {
		return i, err
	}
	tsRepo, tsUser := dis.InstanceGraphDB(slug)

	// generate drupal data
	drPassword, err := dis.NewPassword()
	if err != nil {
		return i, err
	}
	drUser := "admin"

	// make the instance object!
	instance := bookkeeping.Instance{
		Slug: slug,

		OwnerEmail:             "",
		AutoBlindUpdateEnabled: true,

		FilesystemBase: dis.InstanceDir(slug),

		SqlDatabase: sqlDB,
		SqlUser:     sqlUser,
		SqlPassword: sqlPassword,

		GraphDBRepository: tsRepo,
		GraphDBUser:       tsUser,
		GraphDBPassword:   tsPassword,
	}

	i.DrupalUsername = drUser
	i.DrupalPassword = drPassword

	// store the instance in the object and return it!
	i.Instance = instance
	i.dis = dis
	return i, nil
}
