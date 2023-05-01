<?php

// This file contains code to execute a php execution server.
// It is passed as a *command line literal * directly to 'drush:script'.
//
// As such it is preprocessed and shortened.
// This preprocessing is script-specific, and any changes in here might break that optimization.
// It should only contain comments at the beginning of each line, and only starting with '//'.
// See server.go.

// This file runs an infinite REPL-like loop.
// It continously reads from standard input, decodes code to pass to eval(), and send the result back.
//
// Commands are read in DEFLATE base64 encoding (to prevent having to send too much over the wire).
// Results are written to STDOUT back in base64 DEFLATE encoding.
// The results are of the form [$result,$error] - $result being the actual object returned, and $error an error string.

// define a json_encode alias, this saves a single character!
// (we also reuse it in the error string)
$E = 'json_encode';

// prevent STDIN from being buffered
stream_set_read_buffer(STDIN,0);

while(1) {
	// stop outputting errors when executing
	ob_start(null,0,PHP_OUTPUT_HANDLER_CLEANABLE);

	// read the next line (which will be code to execute)
	$b = fgets(STDIN) or exit(0);

	// read the next command to parse
	$b = gzinflate(base64_decode($b));

	// execute the code, and json_encode it
	try {
		$j = $E([eval($b),""]);
	} catch(Throwable $t) {
		$j = $E([null,(string)$t]);
	}

	// if something went wrong, return an error.
	if($j === false) {
		$j = '[null,"' . $E . ' Error"]';
	}

    // and write out the result
	ob_end_clean();

	// deflate and base64_encode
	$j = base64_encode(gzdeflate($j));
	fwrite(STDOUT,"$j\n");
}