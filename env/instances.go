package env

import (
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/internal/fsx"
	"github.com/FAU-CDI/wisski-distillery/internal/stack"
	"github.com/alessio/shellescape"
	"github.com/pkg/errors"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var errNoBookkeeping = exit.Error{
	Message:  "instance %q does not exist in bookkeeping table",
	ExitCode: exit.ExitGeneric,
}

var ErrInstanceNotFound = exit.Error{
	Message:  "instance not found",
	ExitCode: exit.ExitGeneric,
}

var errSQL = exit.Error{
	Message:  "Unknown SQL Error %s",
	ExitCode: exit.ExitGeneric,
}

// Instance returns the instance of the WissKI Distillery with the provided slug
func (dis *Distillery) Instance(slug string) (i Instance, err error) {
	sql := dis.SQL()
	if err := sql.Wait(); err != nil {
		return i, err
	}

	table, err := sql.OpenBookkeeping(false)
	if err != nil {
		return i, err
	}

	// find the instance by slug
	query := table.Where(&bookkeeping.Instance{Slug: slug}).Find(&i.Instance)
	switch {
	case query.Error != nil:
		return i, errSQL.WithMessageF(query.Error)
	case query.RowsAffected == 0:
		return i, ErrInstanceNotFound
	default:
		i.dis = dis
		return i, nil
	}
}

// HasInstance checks if the provided instance exists in the bookeeping table
func (dis *Distillery) HasInstance(slug string) (ok bool, err error) {
	sql := dis.SQL()
	if err := sql.Wait(); err != nil {
		return false, err
	}

	table, err := sql.OpenBookkeeping(false)
	if err != nil {
		return false, err
	}

	query := table.Select("count(*) > 0").Where("slug = ?", slug).Find(&ok)
	if query.Error != nil {
		return false, errSQL.WithMessageF(query.Error)
	}
	return
}

// Instances is like InstancesWith, except that when no slugs are provided, it calls AllInstances.
func (dis *Distillery) Instances(slugs ...string) ([]Instance, error) {
	if len(slugs) == 0 {
		return dis.AllInstances()
	}
	return dis.InstancesWith(slugs...)
}

// AllInstances returns all instances of the WissKI Distillery in consistent order.
//
// There is no guarantee that this order remains identical between different api releases; however subsequent invocations are guaranteed to return the same order.
func (dis *Distillery) AllInstances() ([]Instance, error) {
	return dis.findInstances(true, func(table *gorm.DB) *gorm.DB {
		return table
	})
}

// InstancesWith returns all instances where the slug is in the provided list of names.
// The returned instances are reordered in a consistent order.
func (dis *Distillery) InstancesWith(slugs ...string) ([]Instance, error) {
	return dis.findInstances(true, func(table *gorm.DB) *gorm.DB {
		return table.Where("slug IN ?", slugs)
	})
}

// findInstances finds instance objects based on a query in the bookkeeping table
func (dis *Distillery) findInstances(order bool, query func(table *gorm.DB) *gorm.DB) (instances []Instance, err error) {
	sql := dis.SQL()
	if err := sql.Wait(); err != nil {
		return nil, err
	}

	// open the bookkeeping table
	table, err := sql.OpenBookkeeping(false)
	if err != nil {
		return nil, err
	}

	// prepare a query
	find := table
	if order {
		find = find.Order(clause.OrderByColumn{Column: clause.Column{Name: "slug"}, Desc: false})
	}
	if query != nil {
		find = query(find)
	}

	// fetch bookkeeping instances
	var bks []bookkeeping.Instance
	find = find.Find(&bks)
	if find.Error != nil {
		return nil, errSQL.WithMessageF(find.Error)
	}

	// make proper instances
	instances = make([]Instance, len(bks))
	for i, bk := range bks {
		instances[i].Instance = bk
		instances[i].dis = dis
	}

	return instances, nil
}

// Instance represents a bookkeeping instance
type Instance struct {
	bookkeeping.Instance

	// Credentials for the drupal instance
	DrupalUsername string
	DrupalPassword string

	dis *Distillery
}

// Update updates the bookkeeping table with this instance.
func (instance *Instance) Update() error {
	db, err := instance.dis.SQL().OpenBookkeeping(false)
	if err != nil {
		return err
	}

	// it has never been created => we need to create it in the database
	if instance.Instance.Created.IsZero() {
		return db.Create(&instance.Instance).Error
	}

	// Update based on the primary key!
	return db.Where("pk = ?", instance.Instance.Pk).Updates(&instance.Instance).Error
}

// Delete deletes this instance from the bookkeeping table
func (instance *Instance) Delete() error {
	db, err := instance.dis.SQL().OpenBookkeeping(false)
	if err != nil {
		return err
	}

	// doesn't exist => nothing to delete
	if instance.Instance.Created.IsZero() {
		return nil
	}

	// delete it directly
	return db.Delete(&instance.Instance).Error
}

// Shell executes a shell command inside the
func (instance Instance) Shell(io stream.IOStream, argv ...string) (int, error) {
	return instance.Stack().Exec(io, "barrel", "/user_shell.sh", argv...)
}

// Domain returns the full domain name of this instance
func (instance Instance) Domain() string {
	return fmt.Sprintf("%s.%s", instance.Slug, instance.dis.Config.DefaultDomain)
}

// IfHttps returns value if the distillery has https enabled, the empty string otherwise
// TODO: Fix this to be in a proper place
func (dis *Distillery) IfHttps(value string) string {
	if !dis.HTTPSEnabled() {
		return ""
	}
	return value
}

// URL returns the public URL of this instance
func (instance Instance) URL() *url.URL {
	// setup domain and path
	url := &url.URL{
		Host: instance.Domain(),
		Path: "/",
	}

	// use http or https scheme depending on if the distillery has it enabled
	if instance.dis.HTTPSEnabled() {
		url.Scheme = "https"
	} else {
		url.Scheme = "http"
	}

	return url
}

// Stack represents a stack representing this instance
func (instance Instance) Stack() stack.Installable {
	return stack.Installable{
		Stack: stack.Stack{
			Dir: instance.FilesystemBase,
		},
		ContextResource: filepath.Join("resources", "compose", "barrel"),

		EnvFileResource: filepath.Join("resources", "templates", "docker-env", "barrel"),
		EnvFileContext: map[string]string{
			"REAL_PATH": instance.FilesystemBase,

			"SLUG":         instance.Slug,
			"VIRTUAL_HOST": instance.Domain(),

			"LETSENCRYPT_HOST":  instance.dis.IfHttps(instance.Domain()),
			"LETSENCRYPT_EMAIL": instance.dis.IfHttps(instance.dis.Config.CertbotEmail),

			"UTILS_DIR":                   instance.dis.RuntimeUtilsDir(),
			"GLOBAL_AUTHORIZED_KEYS_FILE": instance.dis.Config.GlobalAuthorizedKeysFile,
		},

		CopyContextFiles: nil,

		TouchFiles: []string{
			"authorized_keys",
		},

		MakeDirsPerm: fs.ModeDir | fs.ModePerm,
		MakeDirs:     []string{"data", ".composer"},
	}
}

func (instance Instance) ReserveStack() stack.Installable {
	return stack.Installable{
		Stack: stack.Stack{
			Dir: instance.FilesystemBase,
		},
		ContextResource: filepath.Join("resources", "compose", "reserve"),

		EnvFileResource: filepath.Join("resources", "templates", "docker-env", "reserve"),
		EnvFileContext: map[string]string{
			"VIRTUAL_HOST": instance.Domain(),

			"LETSENCRYPT_HOST":  instance.dis.IfHttps(instance.Domain()),
			"LETSENCRYPT_EMAIL": instance.dis.IfHttps(instance.dis.Config.CertbotEmail),
		},
	}
}

// Provision provisions an instance, assuming that the required databases already exist.
func (instance Instance) Provision(io stream.IOStream) error {

	// create the basic st!
	st := instance.Stack()
	if err := st.Install(io, stack.InstallationContext{}); err != nil {
		return err
	}

	// Pull and build the stack!
	if err := st.Update(io, false); err != nil {
		return err
	}

	provisionParams := []string{
		instance.Domain(),

		instance.SqlDatabase,
		instance.SqlUser,
		instance.SqlPassword,

		instance.GraphDBRepository,
		instance.GraphDBUser,
		instance.GraphDBPassword,

		instance.DrupalUsername,
		instance.DrupalPassword,

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
		return errors.New("Unable to run provision script")
	}

	return nil
}

func (instance *Instance) NoPrefix() bool {
	return fsx.IsFile(filepath.Join(instance.FilesystemBase, "prefixes.skip"))
}

var errPrefixExecFailed = errors.New("PrefixConfig: Failed to call list_uri_prefixes")

// PrefixConfig returns the prefix config belonging to this instance.
func (instance *Instance) PrefixConfig() (config string, err error) {
	// if the user requested to skip the prefix, then don't do anything with it!
	if instance.NoPrefix() {
		return "", nil
	}

	var builder strings.Builder

	// domain
	builder.WriteString(instance.URL().String() + ":")
	builder.WriteString("\n")

	// default prefixes
	wu := stream.NewIOStream(&builder, nil, nil, 0)
	code, err := instance.Stack().Exec(wu, "barrel", "/bin/bash", "/user_shell.sh", "-c", "drush php:script /wisskiutils/list_uri_prefixes.php")
	if err != nil || code != 0 {
		return "", errPrefixExecFailed
	}

	// custom prefixes
	prefixPath := filepath.Join(instance.FilesystemBase, "prefixes")
	if fsx.IsFile(prefixPath) {
		prefix, err := os.Open(prefixPath)
		if err != nil {
			return "", err
		}
		defer prefix.Close()
		if _, err := io.Copy(&builder, prefix); err != nil {
			return "", err
		}
		builder.WriteString("\n")
	}

	// and done!
	return builder.String(), nil
}
