// Package component holds the main abstraction for components.
//
//spellchecker:words component
package component

//spellchecker:words reflect strconv strings github wisski distillery internal config
import (
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/config"
)

// Components represents a logical subsystem of the distillery.
// A Component should be implemented as a pointer to a struct.
// Every component must embed [Base] and should be initialized using [Init] inside a [lifetime.Lifetime].
//
// By convention these are defined within their corresponding subpackage.
// This subpackage also contains all required resources.
type Component interface {
	// Name returns the name of this component
	// Name should be implemented by the [ComponentBase] struct.
	Name() string

	// ID returns a unique id of this component
	// ID should be implemented by the [ComponentBase] struct.
	ID() string

	// getBase returns the underlying ComponentBase object of this Component.
	// It is used internally during initialization
	getBase() *Base
}

// Base is embedded into every Component
type Base struct {
	name, id string // name and id of this component
	still    Still  // the underlying still of the distillery
}

//lint:ignore U1000 used to implement the private methods of [Component]
func (cb *Base) getBase() *Base {
	return cb
}

// Init initialzes a new componeont Component with the provided still.
// Init is only initended to be used within a lifetime.Lifetime[Component,Still].
func Init(component Component, core Still) {
	base := component.getBase() // pointer to a struct
	base.still = core

	tp := reflect.TypeOf(component).Elem()
	base.name = strings.ToLower(tp.Name())
	base.id = tp.PkgPath() + "." + tp.Name()
}

func (cb Base) Name() string {
	return cb.name
}

func (cb Base) ID() string {
	return cb.id
}

// GetStill returns the still underlying the provided component.
func GetStill(c Component) Still {
	return c.getBase().still
}

// Still represents the central part of a distillery.
// It holds configuration of the distillery.
type Still struct {
	Config   *config.Config // the configuration of the distillery
	Upstream Upstream
}

// Upstream contains the configuration for accessing remote configuration.
type Upstream struct {
	SQL         HostPort
	Triplestore HostPort
	Solr        HostPort
}

func (us Upstream) SQLAddr() string {
	return us.SQL.String()
}

func (us Upstream) TriplestoreAddr() string {
	return us.Triplestore.String()
}

func (us Upstream) SolrAddr() string {
	return us.Solr.String()
}

type HostPort struct {
	Host string
	Port uint32
}

func (hp HostPort) String() string {
	return net.JoinHostPort(hp.Host, strconv.Itoa(int(hp.Port)))
}
