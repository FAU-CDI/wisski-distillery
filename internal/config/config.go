// Package config contains distillery configuration
//
//spellchecker:words config
package config

//spellchecker:words hash math rand reflect time github pkglib reflectx yamlx gopkg yaml embed
import (
	"hash/fnv"
	"math/rand"
	"reflect"
	"time"

	"github.com/tkw1536/pkglib/reflectx"
	"github.com/tkw1536/pkglib/yamlx"
	"gopkg.in/yaml.v3"

	_ "embed"
)

// Config represents the configuration of a WissKI Distillery.
//
// Config is read from a byte stream using [Unmarshal].
//
// Config contains many methods that do not require any interaction with any running components.
// Methods that require running components are instead store inside the [Distillery] or an appropriate [Component].
type Config struct {
	Listen ListenConfig `yaml:"listen" recurse:"true"`
	Paths  PathsConfig  `yaml:"paths" recurse:"true"`
	HTTP   HTTPConfig   `yaml:"http" recurse:"true"`
	Home   HomeConfig   `yaml:"home" recurse:"true"`
	Docker DockerConfig `yaml:"docker" recurse:"true"`

	SQL SQLConfig `yaml:"sql" recurse:"true"`
	TS  TSConfig  `yaml:"triplestore" recurse:"true"`

	// Maximum age for backup in days
	MaxBackupAge time.Duration `yaml:"age" validate:"duration"`

	// Various components use password-based-authentication.
	// These passwords are generated automatically.
	// This variable can be used to determine their length.
	PasswordLength int `yaml:"password_length" default:"64" validate:"positive"`

	// session secret holds the secret for login
	SessionSecret string `yaml:"session_secret" validate:"nonempty" sensitive:"true"`

	// interval to trigger distillery cron tasks in
	CronInterval time.Duration `yaml:"cron_interval" default:"10m" validate:"duration"`

	// ConfigPath is the path this configuration was loaded from (if any)
	ConfigPath string `yaml:"-"`
}

func zeroSensitive(v reflect.Value) {
	reflectx.IterateFields(v.Type(), func(field reflect.StructField, index int) (stop bool) {
		// if we set the recurse tag, recurse into it
		if _, ok := field.Tag.Lookup("recurse"); ok {
			zeroSensitive(v.FieldByName(field.Name))
		}

		// if the field is sensitive, set the zero value!
		if _, ok := field.Tag.Lookup("sensitive"); ok {
			v.FieldByName(field.Name).Set(reflect.Zero(field.Type))
		}
		return false
	})
}

func (config Config) MarshalSensitive() string {
	// zero out all the sensitive fields
	zeroSensitive(reflect.ValueOf(&config).Elem())

	// marshal the result
	result, err := Marshal(&config, nil)
	if err != nil {
		return ""
	}

	return string(result)
}

//go:embed config.yml
var configBytes []byte

// Marshal marshals this configuration in nicely formatted form.
// Where possible, this will maintain yaml comments.
//
// Previous may optionally provide the bytes of a previous configuration file to replace settings in.
// The previous yaml file must be a valid configuration yaml, meaning all fields should be set.
// When previous is of length 0, the default configuration yaml will be used instead.
func Marshal(config *Config, previous []byte) ([]byte, error) {
	if len(previous) == 0 {
		previous = configBytes
	}

	// load the template yaml
	template := new(yaml.Node)
	if err := yaml.Unmarshal(previous, template); err != nil {
		return nil, err
	}

	// load the config yaml
	cfg, err := yamlx.Marshal(config)
	if err != nil {
		return nil, err
	}

	// transplant the configuration yaml into the template
	if err := yamlx.Transplant(template, cfg, true); err != nil {
		return nil, err
	}

	// marshal it again as a set of bytes
	return yaml.Marshal(template)
}

// CSRFSecret return the csrfSecret derived from the session secret
func (config *Config) CSRFSecret() []byte {
	// take the hash of the secret
	h := fnv.New32a()
	h.Write([]byte(config.SessionSecret))

	// seed a random number generator
	rand := rand.New(rand.NewSource(int64(h.Sum32())))

	// take a bunch of bytes from it
	secret := make([]byte, 32)
	rand.Read(secret)
	return secret
}
