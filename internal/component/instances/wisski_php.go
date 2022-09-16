package instances

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/tkw1536/goprogram/stream"
)

var ErrExecInvalidCode = errors.New("invalid code to execute")
var ErrExecNonZero = errors.New("script returned non-zero code")

// ExecPHPScript executes the PHP code as a script within the wisski instance.
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
//
// Standard input and output streams should not be used.
// Standard error is redirected to io.
func (wisski *WissKI) ExecPHPScript(io stream.IOStream, result any, code string, entrypoint string, args ...any) error {
	// make sure the beginning is right
	if !strings.HasPrefix(code, "<?php") {
		return ErrExecInvalidCode
	}

	// make sure that args is not nil, but an array of length 0!
	if args == nil {
		args = []any{}
	}

	// encode code, args and entrypoint!
	codeEscape, err := marshalPHP("?>" + code)
	if err != nil {
		return err
	}

	entrypointEscape, err := marshalPHP(entrypoint)
	if err != nil {
		return err
	}

	argsEscape, err := marshalPHP(args)
	if err != nil {
		return err
	}

	// assemble the script
	script := `
	ob_start(null, 0, PHP_OUTPUT_HANDLER_CLEANABLE);
	eval(` + codeEscape + `);
	ob_end_clean();

	call_user_func(function(){
		ob_start(null, 0, PHP_OUTPUT_HANDLER_CLEANABLE);
		$result = call_user_func_array(` + entrypointEscape + `, ` + argsEscape + `);
		ob_end_clean();
		echo json_encode($result);
	});
`

	// run the script
	var output bytes.Buffer
	res, err := wisski.Shell(io.Streams(&output, nil, strings.NewReader(script), 0), "-c", "drush php:script -")
	if res != 0 {
		return ErrExecNonZero
	}
	if err != nil {
		return err
	}

	// did not request to receive a result
	if result == nil {
		return nil
	}

	// decode the output
	return json.NewDecoder(&output).Decode(result)
}

const marshalRune = 'F' // press to pay respect

// marshalPHP marshals some data which can be marshaled using [json.Encode] into a PHP Expression.
// the string can be safely used directly within php.
func marshalPHP(data any) (string, error) {
	// this function uses json as a data format to transport the data into php.
	// then we build a heredoc to encode it safely, and decode it in php

	// Step 1: Encode the data as json
	jbytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	jstring := string(jbytes)

	// Step 2: Find a delimiter for the heredoc.
	// Step 2a: Find the longest sequence of [marshalRune]s inside the encoded string.
	var current, longest int
	for _, r := range jstring {

		if r == marshalRune {
			current++
		} else {
			current = 0
		}

		if current > longest {
			longest = current
		}
	}
	// Step 2b: Build a string of marshalRune that is one longer!
	delim := strings.Repeat(string(marshalRune), longest+1)

	// Step 3: Assemble the encoded string!
	result := "call_user_func(function(){$x=<<<'" + delim + "'\n" + jstring + "\n" + delim + ";return json_decode(trim($x));})" // press to doubt
	return result, nil
}
