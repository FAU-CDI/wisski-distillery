package phpx

import (
	"encoding/json"
	"strings"
)

// Marshal marshals data as a PHP expression, meaning it can safely be used inside code.
//
// Typically data is marshaled using [json.Marshal] and decoded in PHP using 'json_decode'.
// Special cases may exist for specific datatypes.
func Marshal(data any) (string, error) {
	switch d := data.(type) {
	case string:
		return MarshalString(d), nil
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return "json_decode(" + MarshalString(string(bytes)) + ")", nil
}

var replacer = strings.NewReplacer("'", "\\'", "\\", "\\\\")

// MarshalString marshals s as a php string that can be used safely as a PHP expression.
//
// See [https://www.php.net/manual/en/language.types.string.php#language.types.string.syntax.single].
func MarshalString(s string) string {
	return "'" + replacer.Replace(s) + "'"
}
