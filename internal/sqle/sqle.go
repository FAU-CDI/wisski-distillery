package sqle

import (
	"github.com/feiin/sqlstring"
)

// Format formats the provided query with the given parameters.
func Format(query string, params ...interface{}) string {
	return sqlstring.Format(query, params...)
}
