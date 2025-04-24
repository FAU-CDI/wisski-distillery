package php

//spellchecker:words context strings github wisski distillery internal phpx ingredient barrel
import (
	"context"
	"fmt"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/pkg/errwrap"
)

type PHP struct {
	ingredient.Base
	dependencies struct {
		Barrel *barrel.Barrel
	}
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
func (php *PHP) ExecScript(ctx context.Context, server *phpx.Server, value any, code string, entrypoint string, args ...any) (err error) {
	if server == nil {
		server = php.NewServer()
		defer errwrap.Close(server, "server", &err)
	}

	if code != "" {
		if err := server.MarshalEval(ctx, nil, strings.TrimPrefix(code, "<?php")); err != nil {
			return fmt.Errorf("failed to evaluate code: %w", err)
		}
	}

	if err := server.MarshalCall(ctx, value, entrypoint, args...); err != nil {
		return fmt.Errorf("failed to marshal call: %w", err)
	}
	return nil
}

func (php *PHP) EvalCode(ctx context.Context, server *phpx.Server, value any, code string) (err error) {
	if server == nil {
		server = php.NewServer()
		defer errwrap.Close(server, "server", &err)
	}

	if err := server.MarshalEval(ctx, value, code); err != nil {
		return fmt.Errorf("failed to evaluate code: %w", err)
	}
	return nil
}
