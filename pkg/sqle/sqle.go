package sqle

import (
	"github.com/feiin/sqlstring"
)

// TODO: This is really unsafe and shouldn't be used at all.

// Format formats the provided query with the given parameters.
func Format(query string, params ...interface{}) string {
	return sqlstring.Format(query, params...)
}
