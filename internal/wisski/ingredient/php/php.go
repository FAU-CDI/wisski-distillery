package php

import (
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
)

type PHP struct {
	ingredient.Base

	Barrel *barrel.Barrel
}

// ExecScript executes the PHP code as a script on the given server.
// When server is nil, creates a new server and automatically closes it after execution.
// Calling this function repeatedly with server = nil is inefficient.
//
// The script should define a function called entrypoint, and may define additional functions.
//
// Code must start with "<?php" and may not contain a closing tag.
// Code is expected not to mess with PHPs output buffer.
// Code should not contain user input.
// Code breaking these conventions may or may not result in an error.
//
// It's arguments are encoded as json using [json.Marshal] and decoded within php.
//
// The return value of the function is again marshaled with json and returned to the caller.
func (php *PHP) ExecScript(server *Server, value any, code string, entrypoint string, args ...any) (err error) {
	if server == nil {
		server, err = php.NewServer()
		if err != nil {
			return
		}
		defer server.Close()
	}

	if code != "" {
		if err := server.MarshalEval(nil, strings.TrimPrefix(code, "<?php")); err != nil {
			return err
		}
	}

	return server.MarshalCall(value, entrypoint, args...)
}

func (php *PHP) EvalCode(server *Server, value any, code string) (err error) {
	if server == nil {
		server, err = php.NewServer()
		if err != nil {
			return
		}
		defer server.Close()
	}

	return server.MarshalEval(value, code)
}
