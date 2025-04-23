//spellchecker:words phpx
package phpx

//spellchecker:words encoding json errors strconv time
import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Boolean represents a boolean php value.
//
// The value can be marshaled to and from php and will behave as a PHP would behave.
//
// The value will always be marshaled as "true" or "false".
// Unmarshaling uses [AsBoolean].
type Boolean bool

// AsBoolean tries to cast the given value to a boolean.
//
// It is able to handle any value that would be [json.Unmarshaled] from a corresponding PHP value.
// Value treates all values as the boolean true, except for the ones listed at [doc].
//
// [doc]: https://www.php.net/manual/en/language.types.boolean.php#language.types.boolean.casting
func AsBoolean(value any) (b Boolean, ok bool) {
	switch d := value.(type) {
	case bool:
		return Boolean(d), true
	case float64:
		return d != 0, true
	case string:
		return (d != "" && d != "0"), true
	case []any:
		return len(d) != 0, true
	case map[string]any:
		return len(d) != 0, true
	case nil:
		return true, true
	}
	return true, false
}

func (b Boolean) MarshalJSON() ([]byte, error) {
	if b {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

var errNotABoolean = errors.New("`Boolean': not an integer")

func (b *Boolean) UnmarshalJSON(data []byte) (err error) {
	return UnmarshalIntermediate(b, func(a any) (Boolean, error) {
		b, ok := AsBoolean(a)
		if !ok {
			return Boolean(false), errNotABoolean
		}
		return b, nil
	}, data)
}

// String represents a string php value.
//
// The value can be marshaled to and from php and will behave as a PHP would behave.
//
// The value will always be marshaled as a literal string.
// Unmarshaling uses [AsString].
type String string

// AsString tries to cast the given value to a string.
//
// It is able to handle any value that would be [json.Unmarshaled] from a corresponding PHP value.
// Value casting is described at [doc].
//
// [doc]: https://www.php.net/manual/en/language.types.string.php#language.types.string.casting
func AsString(value any) (s String, ok bool) {
	switch d := value.(type) {
	case bool:
		if d {
			return "1", true
		}
		return "", true
	case float64:
		if d == float64(int64(d)) {
			return String(strconv.FormatInt(int64(d), 10)), true
		}
		// TODO: not sure this is entirely correct
		return String(strconv.FormatFloat(d, 'E', 1, 64)), true
	case string:
		return String(d), true
	case []any, map[string]any:
		return "Array", true
	case nil:
		return "", true
	}

	return "", false
}

func (s String) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(string(s))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}
	return bytes, nil
}

var errNotAString = errors.New("`String': not a string")

func (s *String) UnmarshalJSON(data []byte) (err error) {
	return UnmarshalIntermediate(s, func(a any) (String, error) {
		s, ok := AsString(a)
		if !ok {
			return s, errNotAString
		}
		return s, nil
	}, data)
}

// Integer represents a boolean integer value.
//
// The value can be marshaled to and from php and will behave as a PHP would behave.
//
// The value will always be marshaled as an integer directly
// Unmarshaling uses [AsInteger].
type Integer int64

// AsInteger tries to cast the given value to an integer.
//
// It is able to handle any value that would be [json.Unmarshaled] from a corresponding PHP value.
// Value casting is described at [doc].
//
// [doc]: https://www.php.net/manual/en/language.types.integer.php#language.types.integer.casting
func AsInteger(value any) (i Integer, ok bool) {
	str, ok := AsString(value)
	if !ok {
		return 0, false
	}

	// try to parse the "leading" string, by successively cutting off parts of the tail
	// once we have a valid number, return it.
	for l := range len(str) {
		i64, err := strconv.ParseInt(string(str)[:len(str)-l], 10, 64)
		if err != nil {
			continue
		}
		return Integer(i64), true
	}
	return 0, true
}

func (i Integer) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(int64(i))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}
	return bytes, nil
}

var errNotAnInteger = errors.New("`Integer': not an integer")

func (i *Integer) UnmarshalJSON(data []byte) (err error) {
	return UnmarshalIntermediate(i, func(a any) (Integer, error) {
		i, ok := AsInteger(a)
		if !ok {
			return i, errNotAnInteger
		}
		return i, nil
	}, data)
}

// Timestamp represents a time value in PHP, represented as an integer.
//
//nolint:recvcheck
type Timestamp time.Time

func (ts Timestamp) Time() time.Time {
	return time.Time(ts)
}
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(ts.Time().Unix(), 10)), nil
}

func (ts *Timestamp) UnmarshalJSON(data []byte) (err error) {
	return UnmarshalIntermediate(ts, func(value Integer) (Timestamp, error) {
		return Timestamp(time.Unix(int64(value), 0)), nil
	}, data)
}

// UnmarshalIntermediate unmarshals src into dest using an intermediate value of type I.
//
// It first unmarshals src into a new value of type I.
// It then calls parser to parse I into T.
func UnmarshalIntermediate[I, T any](dest *T, parser func(I) (T, error), src []byte) (err error) {
	var temp I
	err = json.Unmarshal(src, &temp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal into json: %w", err)
	}

	*dest, err = parser(temp)
	if err != nil {
		return fmt.Errorf("parser returned error: %w", err)
	}

	return nil
}
