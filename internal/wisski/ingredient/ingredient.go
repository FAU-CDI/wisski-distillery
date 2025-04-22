//spellchecker:words ingredient
package ingredient

//spellchecker:words reflect strings github wisski distillery internal component liquid
import (
	"reflect"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/liquid"
)

// Ingredients represent a part of a WissKI instance.
// An Ingredient should be implemented as a pointer to a struct.
// Every ingredient must embed [Base] and should be initialized using [Init] inside a [lifetime.Lifetime].
//
// By convention these are defined within their corresponding subpackage.
// This subpackage also contains all required resources.
type Ingredient interface {
	// Name returns the name of this ingredient
	// Name should be implemented by the [Base] struct.
	Name() string

	// getBase returns the underlying Base object of this Ingredient.
	// It is used internally during initialization
	getBase() *Base
}

// Base is embedded into every Ingredient.
type Base struct {
	name   string         // name is the name of this ingredient
	liquid *liquid.Liquid // the underlying liquid
}

// GetLiquid gets the liquid of this Ingredient.
func GetLiquid(i Ingredient) *liquid.Liquid {
	return i.getBase().liquid
}

// GetStill returns the still of the distillery associated with the provided ingredient.
func GetStill(i Ingredient) component.Still {
	return component.GetStill(GetLiquid(i).Malt)
}

//lint:ignore U1000 used to implement the private methods of [Component]
func (cb *Base) getBase() *Base {
	return cb
}

// Init initializes a new Ingredient.
// Init is only intended to be used within a lifetime.Lifetime[Ingredient,*Liquid].
func Init(ingredient Ingredient, liquid *liquid.Liquid) {
	base := ingredient.getBase() // pointer to a struct
	base.liquid = liquid
	base.name = strings.ToLower(reflect.TypeOf(ingredient).Elem().Name())
}

func (cb Base) Name() string {
	return cb.name
}
