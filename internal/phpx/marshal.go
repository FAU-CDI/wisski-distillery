//spellchecker:words phpx
package phpx

//spellchecker:words encoding json math strconv strings github pkglib collection
import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/tkw1536/pkglib/collection"
)

// Marshal marshals data as a PHP expression, so that it can be safely used inside a php expession.
//
// Typically data is marshaled using [json.Marshal] and decoded in PHP using 'json_decode'.
// Special cases may exist for specific datatypes.
func Marshal(data any) (string, error) {
	if str, ok := data.(string); ok {
		return MarshalString(str), nil
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	return "json_decode(" + MarshalString(string(bytes)) + ")", nil
}

// No valid value ever returns the empty string.
func MarshalJSON(v any) string {
	switch v := v.(type) {
	case nil:
		return PHPNil
	case bool:
		return MarshalBool(v)
	case float64:
		return MarshalFloat(v)
	case string:
		return MarshalString(v)
	case []any:
		return MarshalSlice(v)
	case map[string]any:
		return MarshalMap(v)
	}
	return ""
}

const (
	// PHPNil represents the equivalent of a nil value in php.
	PHPNil = "null"

	phpTrue  = "!0"
	phpFalse = "!1"

	phpNaN              = "NAN"
	phpPositiveInfinity = "INF"
	phpNegativeInfinity = "-INF"
)

// MarshalBool marshals b as a boolean to be used in php code.
// This corresponds to the strings "true" or "false".
func MarshalBool(b bool) string {
	if b {
		return phpTrue
	}
	return phpFalse
}

// MarshalFloat marshals a floating point number or integer.
func MarshalFloat(f float64) string {
	// if we actually have an integer, return it!
	if i := int64(f); f == float64(i) {
		return MarshalInt(i)
	}

	// special cases
	if math.IsNaN(f) {
		return phpNaN
	}
	if math.IsInf(f, 1) {
		return phpPositiveInfinity
	}
	if math.IsInf(f, -1) {
		return phpNegativeInfinity
	}

	// all other cases
	return strconv.FormatFloat(f, 'E', -1, 64)
}

// MarshalInt marshals an integer as a string to be used inside a php literal.
func MarshalInt(i int64) string {
	return strconv.FormatInt(i, 10)
}

var stringReplacer = strings.NewReplacer("'", "\\'", "\\", "\\\\")

// MarshalString marshals s as a php string that can be used safely as a PHP expression.
func MarshalString(s string) string {
	// See [https://www.php.net/manual/en/language.types.string.php#language.types.string.syntax.single]
	return "'" + stringReplacer.Replace(s) + "'"
}

func MarshalSlice(slice []any) string {
	var builder strings.Builder

	builder.WriteRune('[')
	{
		for _, v := range slice {
			builder.WriteString(MarshalJSON(v))
			builder.WriteRune(',')
		}
	}
	builder.WriteRune(']')

	return builder.String()
}

func MarshalMap(m map[string]any) string {
	var builder strings.Builder

	builder.WriteString("array(")
	for k, v := range collection.IterSorted(m) {
		builder.WriteString(MarshalString(k))
		builder.WriteString("=>")
		builder.WriteString(MarshalJSON(v))
		builder.WriteString(",")
	}
	builder.WriteString(")")

	return builder.String()
}
