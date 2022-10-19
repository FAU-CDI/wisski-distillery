<?php

// This file contains code to execute a php execution server.
// It is passed as a *command line literal * directly to 'drush:script'.
//
// As such it is preprocessed and shortened.
// It should only contain comments at the beginning of each line, and only starting with '//'.
// See server.go.

// prevent STDIN from being buffered
stream_set_read_buffer(STDIN,0);

// stop outputting errors when executing
ob_start(null,0,PHP_OUTPUT_HANDLER_CLEANABLE);


while(1){
	// read the next line to get an end-of-line marker
	$marker = fgets(STDIN);
	if (!$marker) break;

	// accumulate the buffer until the marker is reached
	// bail out if there is an unexpected end of input
	$buffer = "";
	while(1) {
		$line = fgets(STDIN);
		if (!$line) break 2;
		if ($line === $marker) break;
		$buffer .= $line . "\n";
	}

    // execute it
	try{
		$json = json_encode([eval($buffer),""]);
	}catch(Throwable $t){
		$json = json_encode([null,(string)$t]);
	}
	if($json===false) {
		$json = '[null,"Error encoding result"]';
	}

    // and write out the result
	ob_end_clean();
	fwrite(STDOUT,"$json\n");
    ob_start(null,0,PHP_OUTPUT_HANDLER_CLEANABLE);
}