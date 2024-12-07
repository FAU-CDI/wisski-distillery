//spellchecker:words tokens
package tokens

//spellchecker:words errors http strings github wisski distillery internal component models golang slices
import (
	"errors"
	"net/http"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"golang.org/x/exp/slices"
)

const (
	authHeader = "Authorization" // authorization
	authBearer = "Bearer" + " "  // Prefix for bearer
)

// TokenOf returns the token header found in the given request.
// If r is nil, or there is no token, returns nil.
// Error is only set if there is an error accessing the table that stores tokens.
func (tok *Tokens) TokenOf(r *http.Request) (*models.Token, error) {
	if r == nil {
		return nil, nil
	}

	// make sure that the authorization header exists and starts with the bearer
	auth := r.Header.Get(authHeader)
	if !strings.HasPrefix(auth, authBearer) {
		return nil, nil
	}

	// get the token
	id := strings.TrimSpace(strings.TrimPrefix(auth, authBearer))
	if id == "" {
		return nil, nil
	}

	table, err := tok.table(r.Context())
	if err != nil {
		return nil, err
	}

	// take a single object from the tokens
	var tokenObj models.Token
	res := table.Where(&models.Token{Token: id}).Find(&tokenObj)

	if res.Error != nil {
		return nil, errors.Join(ErrNoToken, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}

	// and return the token object
	return &tokenObj, nil
}

var ErrNoToken = errors.New("no token")

// Check checks if there is a token in the given request and if this request has an appropriate token with the appropriate scope.
//
// If the token is found and has the requested token, returns true, nil.
// If there is a token found, but the specific scope is not set, returns false, nil.
// If there is no valid authentication token found, returns false and an error that wraps ErrNoToken.
// In other cases, other errors may be returned.
//
// Note that the scope may require an parameter to be validated.
// This validation should take place in the appropriate ScopeProvider; which should recursively invoke this method.
func (tok *Tokens) Check(r *http.Request, scope component.Scope) (bool, error) {
	// get the token object from the request
	tokenObj, err := tok.TokenOf(r)
	if tokenObj == nil {
		if err == nil {
			return false, ErrNoToken
		}
		return false, errors.Join(ErrNoToken, err)
	}

	// TODO: Do we need this function?

	// get the scopes
	scopes := tokenObj.GetScopes()
	if scopes == nil {
		// all scopes (implicitly)
		return true, nil
	}

	// else check if they are contained
	return slices.Contains(scopes, string(scope)), nil
}
