<?php

// This file contains code to execute a php execution server.
// It is passed as a *command line literal * directly to 'drush:script'.
//
// As such it is preprocessed and shortened.
// It should only contain comments at the beginning of each line, and only starting with '//'.
// See wisski_php_server.go.

// don't buffer stdin!
stream_set_read_buffer(STDIN,0);

// stop all other output
ob_start(null,0,PHP_OUTPUT_HANDLER_CLEANABLE);

while($line =  fgets(STDIN)){
    // decode the command to run
	$code=@json_decode($line);

    // execute it
	try{
		$json = json_encode([eval($code),""]);
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