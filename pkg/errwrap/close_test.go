package errwrap_test

import (
	"errors"
	"testing"

	"github.com/FAU-CDI/wisski-distillery/pkg/errwrap"
)

var (
	errCloseNotOK    = errors.New("close not ok")
	errFunctionNotOK = errors.New("function not ok")
)

// CloseIfOK is an [io.Closer] that returns an error iff it is false.
type CloseIfOK bool

func (ok CloseIfOK) Close() error {
	if !ok {
		return errCloseNotOK
	}
	return nil
}

func TestClose(t *testing.T) {
	t.Parallel()

	someFunction := func(functionOK, closeOK bool) (e error) {
		defer errwrap.Close(CloseIfOK(closeOK), "object", &e)

		// return if everything is fine
		if !functionOK {
			return errFunctionNotOK
		}
		return nil
	}

	for _, tt := range []struct {
		name                string
		functionOK, closeOK bool

		wantError bool  // did we want an error at all?
		wantExact error // did we want to be exactly this error? (only checked if != nil)

		wantCloseNotOK    bool // did we want to errors.Is(errCloseNotOK)
		wantFunctionNotOK bool // did we want to errors.Is(errFunctionNotOK)
	}{
		{
			name:       "normal executation",
			functionOK: true, closeOK: true,

			wantError:         false,
			wantExact:         nil,
			wantCloseNotOK:    false,
			wantFunctionNotOK: false,
		},

		{
			name:       "only function fails",
			functionOK: false, closeOK: true,

			wantError:         true,
			wantExact:         errFunctionNotOK,
			wantCloseNotOK:    false,
			wantFunctionNotOK: true,
		},

		{
			name:       "only close fails",
			functionOK: true, closeOK: false,

			wantError:         true,
			wantCloseNotOK:    true,
			wantFunctionNotOK: false,
		},

		{
			name:       "both fail",
			functionOK: false, closeOK: false,

			wantError:         true,
			wantCloseNotOK:    true,
			wantFunctionNotOK: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := someFunction(tt.functionOK, tt.closeOK)

			gotError := (err != nil)
			if gotError != tt.wantError {
				t.Errorf("got error = %v, but wanted error = %v", gotError, tt.wantError)
			}
			gotExact := (tt.wantExact == nil) || (err == tt.wantExact)
			if !gotExact {
				t.Errorf("got error = %v, but wanted error = %v", err, tt.wantExact)
			}

			gotCloseNotOK := errors.Is(err, errCloseNotOK)
			if gotCloseNotOK != tt.wantCloseNotOK {
				t.Errorf("got errors.Is(err, errCloseNotOK) = %v, but wanted = %v", gotCloseNotOK, tt.wantCloseNotOK)
			}

			gotFunctionNotOK := errors.Is(err, errFunctionNotOK)
			if gotFunctionNotOK != tt.wantFunctionNotOK {
				t.Errorf("got errors.Is(err, errFunctionNotOK) = %v, but wanted = %v", gotFunctionNotOK, tt.wantFunctionNotOK)
			}
		})
	}
}
