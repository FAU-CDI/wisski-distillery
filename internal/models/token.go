package models

import (
	"encoding/json"
)

// TokensTable is the name of the table the 'Token' model is stored in.
const TokensTable = "tokens"

// Token represents an access token for a specific user
type Token struct {
	Pk uint `gorm:"column:pk;primaryKey"`

	Token string `gorm:"column:token;unique:true;not null"`
	User  string `gorm:"column:user;not null"` // (distillery) username

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
