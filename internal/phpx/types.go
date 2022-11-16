package phpx

import (
	"encoding/json"
	"strconv"
	"time"
)

// PHPBoolean represents a boolean php value.
//
// The value can be marshaled to and from php and will behave as a PHP would behave.
//
// The value will always be marshaled as "true" or "false".
// Unmarshaling uses [Boolean].
type PHPBoolean bool

func (bi PHPBoolean) MarshalJSON() ([]byte, error) {
	if bi {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}
func (bi *PHPBoolean) UnmarshalJSON(data []byte) (err error) {
	// unmarshal into a generic value
	var value any
	err = json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	// cast into a boolean
	cast, ok := Boolean(value)
	if !ok {
		value = false
	}
	*bi = PHPBoolean(cast)

	return nil
}

// Boolean tries to cast the given value to a boolean.
//
// It is able to handle any value that would be [json.Unmarshaled] from a corresponding PHP value.
// Value treates all values as the boolean true, except for the ones listed at [doc].
//
// [doc]: https://www.php.net/manual/en/language.types.boolean.php#language.types.boolean.casting
func Boolean(value any) (b bool, ok bool) {
	switch d := value.(type) {
	case bool:
		return d, true
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

// String tries to cast the given value to a string.
//
// It is able to handle any value that would be [json.Unmarshaled] from a corresponding PHP value.
// Value casting is described at [doc].
//
// [doc]: https://www.php.net/manual/en/language.types.string.php#language.types.string.casting
func String(value any) (s string, ok bool) {
	switch d := value.(type) {
	case bool:
		if d {
			return "1", true
		}
		return "", true
	case float64:
		if d == float64(int64(d)) {
			return strconv.FormatInt(int64(d), 10), true
		}
		// TODO: not sure this is entirely correct
		// and we should handle ints here!
		return strconv.FormatFloat(d, 'E', 1, 64), true
	case string:
		return d, true
	case []any, map[string]any:
		return "Array", true
	case nil:
		return "", true
	}

	return "", false
}

// Integer tries to cast the given value to an integer.
//
// It is able to handle any value that would be [json.Unmarshaled] from a corresponding PHP value.
// Value casting is described at [doc].
//
// [doc]: https://www.php.net/manual/en/language.types.integer.php#language.types.integer.casting
func Integer(value any) (i int64, ok bool) {
	str, ok := String(value)
	if !ok {
		return 0, false
	}

	// try to parse the "leading" string, by successively cutting off parts of the tail
	// once we have a valid number, return it.
	for l := 0; l < len(str); l++ {
		i64, err := strconv.ParseInt(str[:len(str)-l], 10, 64)
		if err != nil {
			continue
		}
		return i64, true
	}
	return 0, true
}

// TimeInt represents a time value in PHP, represented as an integer
type TimeInt time.Time

func (ts TimeInt) Time() time.Time {
	return time.Time(ts)
}
func (ts TimeInt) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(ts.Time().Unix(), 10)), nil
}

func (ts *TimeInt) UnmarshalJSON(data []byte) (err error) {
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	unix, _ := Integer(value)
	*ts = TimeInt(time.Unix(unix, 0))
	return nil
}
