// Package config contains distillery configuration
//
//spellchecker:words config
package config

//spellchecker:words hash math rand reflect time pkglib reflectx yamlx golang crypto scrypt gopkg yaml embed
import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"reflect"
	"time"

	"go.tkw01536.de/pkglib/reflectx"
	"go.tkw01536.de/pkglib/yamlx"
	"golang.org/x/crypto/scrypt"
	"gopkg.in/yaml.v3"

	_ "embed"
)

// Config represents the configuration of a WissKI Distillery.
//
// Config is read from a byte stream using [Unmarshal].
//
// Config contains many methods that do not require any interaction with any running components.
// Methods that require running components are instead store inside the [Distillery] or an appropriate [Component].
//
//nolint:recvcheck
type Config struct {
	Listen ListenConfig `recurse:"true" yaml:"listen"`
	Paths  PathsConfig  `recurse:"true" yaml:"paths"`
	HTTP   HTTPConfig   `recurse:"true" yaml:"http"`
	Home   HomeConfig   `recurse:"true" yaml:"home"`
	Docker DockerConfig `recurse:"true" yaml:"docker"`

	SQL SQLConfig `recurse:"true" yaml:"sql"`
	TS  TSConfig  `recurse:"true" yaml:"triplestore"`

	// Maximum age for backup in days
	MaxBackupAge time.Duration `validate:"duration" yaml:"age"`

	// Various components use password-based-authentication.
	// These passwords are generated automatically.
	// This variable can be used to determine their length.
	PasswordLength int `default:"64" validate:"positive" yaml:"password_length"`

	// session secret holds the secret for login
	SessionSecret string `sensitive:"true" validate:"nonempty" yaml:"session_secret"`

	// interval to trigger distillery cron tasks in
	CronInterval time.Duration `default:"10m" validate:"duration" yaml:"cron_interval"`

	// ConfigPath is the path this configuration was loaded from (if any)
	ConfigPath string `yaml:"-"`
}

func zeroSensitive(v reflect.Value) {
	for field := range reflectx.IterFields(v.Type()) {
		// if we set the recurse tag, recurse into it
		if _, ok := field.Tag.Lookup("recurse"); ok {
			zeroSensitive(v.FieldByName(field.Name))
		}

		// if the field is sensitive, set the zero value!
		if _, ok := field.Tag.Lookup("sensitive"); ok {
			v.FieldByName(field.Name).Set(reflect.Zero(field.Type))
		}
	}
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
		return nil, fmt.Errorf("failed to unmarshal previous configuration: %w", err)
	}

	// load the config yaml
	cfg, err := yamlx.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// transplant the configuration yaml into the template
	if err := yamlx.Transplant(template, cfg, true); err != nil {
		return nil, fmt.Errorf("failed to render template with configuration: %w", err)
	}

	// marshal it again as a set of bytes
	out, err := yaml.Marshal(template)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}
	return out, nil
}

// SessionKey returns a key used for sessions to be derived from the session secret.
func (config *Config) SessionKey() []byte {
	return config.derivedKey(32)
}

// CSRFSecret return the csrfSecret derived from the session secret.
func (config *Config) CSRFKey() []byte {
	return config.derivedKey(0)
}

// deriveKey derives a 32-bit key which can be used with sessions.
// It will change when the config file changes.
func (config *Config) derivedKey(skip int) []byte {
	salt := config.makeSalt(skip, 64)

	bytes, err := scrypt.Key([]byte(config.SessionSecret), salt, 32768, 8, 1, 32)
	if err != nil {
		panic(fmt.Sprintf("scrypt: derivedKey returned error: %s", err))
	}

	return bytes
}

// makeSalt makes some salt for key deriviation.
// It is based on the contents of the config file.
func (config *Config) makeSalt(skip, size int) []byte {
	// TODO: Generate random salt and read it from somewhere!
	h := fnv.New64a()
	if _, err := h.Write([]byte(config.MarshalSensitive())); err != nil {
		panic("hash failed to write")
	}
	sum := int64(h.Sum64()) // #nosec G115 -- this wraps around, but that's fine!

	// initialize the PRNG and go forward
	rand := rand.New(rand.NewSource(sum)) // #nosec G404 -- this is used to make salt only
	for range skip {
		rand.Int63()
	}

	// and get the bytes
	salt := make([]byte, size)
	if _, err := rand.Read(salt); err != nil {
		panic("never reached: rand.Read() always returns err == nil")
	}
	return salt
}
