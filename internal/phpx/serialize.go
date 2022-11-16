package phpx

import "encoding/json"

// BooleanIsh represents a boolean php value.
//
// The value can be serialized to and from php and will behave accordingly.
//
// The value will always be Marshaled as "true" or "false".
//
// When Unmarshaled, it behaves as described on https://www.php.net/manual/en/language.types.boolean.php#language.types.boolean.casting.
type BooleanIsh bool

func (bi BooleanIsh) MarshalJSON() ([]byte, error) {
	if bi {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}
func (bi *BooleanIsh) UnmarshalJSON(data []byte) (err error) {
	// unmarshal into a generic value
	var value any
	err = json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	// check if it is false ish
	var isFalseIsh bool
	switch d := value.(type) {
	case bool:
		isFalseIsh = !d
	case int:
		isFalseIsh = d == 0
	case float64:
		isFalseIsh = d == 0
	case string:
		isFalseIsh = d == "" || d == "0"
	case []any:
		isFalseIsh = len(d) == 0
	case map[string]any:
		isFalseIsh = len(d) == 0
	case nil:
		isFalseIsh = true
	}
	*bi = BooleanIsh(!isFalseIsh)

	return nil
}
