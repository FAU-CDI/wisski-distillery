//spellchecker:words models
package models

//spellchecker:words encoding json
import (
	"encoding/json"
)

// TokensTable is the name of the table the 'Token' model is stored in.
const TokensTable = "tokens"

// Token represents an access token for a specific user
type Token struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	Token   string `json:"-" gorm:"column:token;unique:true;not null"` // token used by the actual api (shown only once)
	TokenID string `gorm:"column:id;unique:true;not null"`             // token id (displayed to user, used for finding it)

	User string `gorm:"column:user;not null"` // (distillery) username

	Description string `gorm:"column:description"`

	AllScopes bool   `gorm:"column:all;not null"`
	Scopes    []byte `gorm:"column:scopes;not null"` // comma-seperated list of scopes
}

// GetScopes returns the scopes associated with this Token.
//
// If this token implicitly has all scopes, returns nil.
// If this token has no scopes, returns an empty string slice.
func (token *Token) GetScopes() (scopes []string) {
	// all scopes
	if token.AllScopes {
		return nil
	}

	// unmarshal the scopes associated with this token
	// and ensure that it is never nil.
	err := json.Unmarshal(token.Scopes, &scopes)
	if scopes == nil || err != nil {
		scopes = []string{}
	}
	return
}

// SetScopes sets the scopes associated to this token to scopes.
// It scopes is nil, sets the token to permit all scopes.
func (token *Token) SetScopes(scopes []string) {
	token.AllScopes = scopes == nil
	if token.AllScopes {
		scopes = []string{}
	}
	token.Scopes, _ = json.Marshal(scopes)
}
