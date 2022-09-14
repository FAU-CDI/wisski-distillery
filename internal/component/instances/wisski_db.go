package instances

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/pkg/fsx"
	"github.com/alessio/shellescape"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// WissKI represents a single WissKI Instance
type WissKI struct {
	// Whatever is stored inside the bookkeeping database
	bookkeeping.Instance

	// Credentials to Drupal
	DrupalUsername string
	DrupalPassword string

	// reference to the component!
	instances *Instances
}

// Save saves this instance in the bookkeeping table
func (wisski *WissKI) Save() error {
	db, err := wisski.instances.SQL.OpenBookkeeping(false)
	if err != nil {
		return err
	}

	// it has never been created => we need to create it in the database
	if wisski.Instance.Created.IsZero() {
		return db.Create(&wisski.Instance).Error
	}

	// Update based on the primary key!
	return db.Where("pk = ?", wisski.Instance.Pk).Updates(&wisski.Instance).Error
}

// Delete deletes this instance from the bookkeeping table
func (wisski *WissKI) Delete() error {
	db, err := wisski.instances.SQL.OpenBookkeeping(false)
	if err != nil {
		return err
	}

	// doesn't exist => nothing to delete
	if wisski.Instance.Created.IsZero() {
		return nil
	}

	// delete it directly
	return db.Delete(&wisski.Instance).Error
}

// Shell executes a shell command inside the
func (wisski WissKI) Shell(io stream.IOStream, argv ...string) (int, error) {
	return wisski.Stack().Exec(io, "barrel", "/bin/sh", append([]string{"/user_shell.sh"}, argv...)...)
}

// Domain returns the full domain name of this instance
func (wisski WissKI) Domain() string {
	return fmt.Sprintf("%s.%s", wisski.Slug, wisski.instances.Config.DefaultDomain)
}

// URL returns the public URL of this instance
func (wisski WissKI) URL() *url.URL {
	// setup domain and path
	url := &url.URL{
		Host: wisski.Domain(),
		Path: "/",
	}

	// use http or https scheme depending on if the distillery has it enabled
	if wisski.instances.Config.HTTPSEnabled() {
		url.Scheme = "https"
	} else {
		url.Scheme = "http"
	}

	return url
}

//go:embed all:instances/barrel instances/barrel.env
var barrelResources embed.FS

// Stack represents a stack representing this instance
func (wisski WissKI) Stack() component.Installable {
	return component.Installable{
		Stack: component.Stack{
			Dir: wisski.FilesystemBase,
		},

		Resources:   barrelResources,
		ContextPath: filepath.Join("instances", "barrel"),
		EnvPath:     filepath.Join("instances", "barrel.env"),

		EnvContext: map[string]string{
			"DATA_PATH": filepath.Join(wisski.FilesystemBase, "data"),

			"SLUG":         wisski.Slug,
			"VIRTUAL_HOST": wisski.Domain(),

			"LETSENCRYPT_HOST":  wisski.instances.Config.IfHttps(wisski.Domain()),
			"LETSENCRYPT_EMAIL": wisski.instances.Config.IfHttps(wisski.instances.Config.CertbotEmail),

			"RUNTIME_DIR":                 wisski.instances.Config.RuntimeDir(),
			"GLOBAL_AUTHORIZED_KEYS_FILE": wisski.instances.Config.GlobalAuthorizedKeysFile,
		},

		MakeDirsPerm: fs.ModeDir | fs.ModePerm,
		MakeDirs:     []string{"data", ".composer"},

		TouchFiles: []string{
			filepath.Join("data", "authorized_keys"),
		},
	}
}

//go:embed all:instances/reserve instances/reserve.env
var reserveResources embed.FS

func (wisski WissKI) ReserveStack() component.Installable {
	return component.Installable{
		Stack: component.Stack{
			Dir: wisski.FilesystemBase,
		},

		Resources:   reserveResources,
		ContextPath: filepath.Join("instances", "reserve"),
		EnvPath:     filepath.Join("instances", "reserve.env"),

		EnvContext: map[string]string{
			"VIRTUAL_HOST": wisski.Domain(),

			"LETSENCRYPT_HOST":  wisski.instances.Config.IfHttps(wisski.Domain()),
			"LETSENCRYPT_EMAIL": wisski.instances.Config.IfHttps(wisski.instances.Config.CertbotEmail),
		},
	}
}

// Provision provisions an instance, assuming that the required databases already exist.
func (wisski WissKI) Provision(io stream.IOStream) error {

	// create the basic st!
	st := wisski.Stack()
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

// NoPrefix checks if this WissKI instance is excluded from generating prefixes
func (wisski *WissKI) NoPrefix() bool {
	return fsx.IsFile(filepath.Join(wisski.FilesystemBase, "prefixes.skip"))
}

var errPrefixExecFailed = errors.New("PrefixConfig: Failed to call list_uri_prefixes")

// PrefixConfig returns the prefix config belonging to this instance.
func (wisski *WissKI) PrefixConfig() (config string, err error) {
	// if the user requested to skip the prefix, then don't do anything with it!
	if wisski.NoPrefix() {
		return "", nil
	}

	var builder strings.Builder

	// domain
	builder.WriteString(wisski.URL().String() + ":")
	builder.WriteString("\n")

	// default prefixes
	wu := stream.NewIOStream(&builder, nil, nil, 0)
	code, err := wisski.Stack().Exec(wu, "barrel", "/bin/bash", "/user_shell.sh", "-c", "drush php:script /wisskiutils/list_uri_prefixes.php")
	if err != nil || code != 0 {
		return "", errPrefixExecFailed
	}

	// custom prefixes
	prefixPath := filepath.Join(wisski.FilesystemBase, "prefixes")
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

var errPathbuildersExecFailed = errors.New("ExportPathbuilders: Failed to call export_pathbuilder")

// ExportPathbuilders writes pathbuilders into the directory dest
func (wisski *WissKI) ExportPathbuilders(dest string) error {
	// export all the pathbuilders into the buffer
	var buffer bytes.Buffer
	wu := stream.NewIOStream(&buffer, nil, nil, 0)
	code, err := wisski.Stack().Exec(wu, "barrel", "/bin/bash", "/user_shell.sh", "-c", "drush php:script /wisskiutils/export_pathbuilder.php")
	if err != nil || code != 0 {
		return errPathbuildersExecFailed
	}

	// decode them as a json array
	var pathbuilders map[string]string
	if err := json.NewDecoder(&buffer).Decode(&pathbuilders); err != nil {
		return err
	}

	// sort the names of the pathbuilders
	names := maps.Keys(pathbuilders)
	slices.Sort(names)

	// write each into a file!
	for _, name := range names {
		pbxml := []byte(pathbuilders[name])
		name := filepath.Join(dest, fmt.Sprintf("%s.xml", name))
		if err := os.WriteFile(name, pbxml, fs.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
