package php

//spellchecker:words context strings github wisski distillery internal phpx ingredient barrel
import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
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
		if err != nil {
			return
		}
		defer func() {
			e2 := server.Close()
			if e2 == nil {
				return
			}
			e2 = fmt.Errorf("failed to close server: %w", e2)
			if err == nil {
				err = e2
			} else {
				err = errors.Join(err, e2)
			}
		}()
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
		defer func() {
			e2 := server.Close()
			if e2 == nil {
				return
			}
			e2 = fmt.Errorf("failed to close server: %w", e2)
			if err == nil {
				err = e2
			} else {
				err = errors.Join(err, e2)
			}
		}()
	}

	if err := server.MarshalEval(ctx, value, code); err != nil {
		return fmt.Errorf("failed to evaluate code: %w", err)
	}
	return nil
}
