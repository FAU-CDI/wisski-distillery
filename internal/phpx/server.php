<?php

// This file contains code to execute a php execution server.
// It is passed as a *command line literal * directly to 'drush:script'.
//
// As such it is preprocessed and shortened.
// This preprocessing is script-specific, and any changes in here might break that optimization.
// It should only contain comments at the beginning of each line, and only starting with '//'.
// See server.go.

// define a json_encode alias, this saves a single character!
// (we also reuse it in the error string)
$E = 'json_encode';

// prevent STDIN from being buffered
stream_set_read_buffer(STDIN,0);

while(1) {
	// stop outputting errors when executing
	ob_start(null,0,PHP_OUTPUT_HANDLER_CLEANABLE);

	// read the next line to get an end-of-line marker
	$m = fgets(STDIN) or exit(0);

	// accumulate the buffer until the marker is reached
	// bail out if there is an unexpected end of input
	for($b = $l = ""; $l !== $m;) {
		$b .= $l;
		$l = fgets(STDIN) or exit(1);
	}

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
	fwrite(STDOUT,"$j\n");
}