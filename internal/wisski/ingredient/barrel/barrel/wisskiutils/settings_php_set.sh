#!/bin/bash

# settings_php_set.sh name value
# Sets the 'settings.php' setting 'name' to 'value'.
# Value must be json-encoded. 

NAME=$1
VALUE=$2

if [ -z "$NAME" ]; then
    echo "Usage: settings_php_set.sh NAME VALUE"
    exit 1
fi;

if [ -z "$VALUE" ]; then
    echo "Usage: settings_php_set.sh NAME VALUE"
    exit 1
fi;

cd /var/www/data/project
chmod u+w web/sites/default/settings.php

(echo "$NAME"; echo "$VALUE" ) | drush php:eval '
    if(is_file(DRUPAL_ROOT . "/internal/")) {
        include_once DRUPAL_ROOT . "/internal/core/includes/install.inc";
    } else {
        include_once DRUPAL_ROOT . "/core/includes/install.inc";
    }

    // read NAME and VALUE from STDIN
    $content=file_get_contents("php://stdin");  
    $newline=strpos($content, "\n");
    $name=trim(substr($content, 0, $newline));
    $jvalue=trim(substr($content, $newline + 1));

    // decode json values
    $value = @json_decode($jvalue);
    if ($value === null && json_last_error() !== JSON_ERROR_NONE) {
        echo "Invalid JSON, cannot update settings.php. \n";
        return 1;
    }

    // make parameters to drush_rewrite_settings
    $settings["settings"][$name] = (object)[
        "value" => $value,
        "required" => TRUE,
    ];

    // find the actual settings.php file to rewrite
    $filename = DRUPAL_ROOT . "/" . \Drupal::service("site.path") . "/settings.php";
    drupal_rewrite_settings($settings, $filename);

    echo "Wrote " . $filename . "\n";
    return 0;   
';
EXIT=$?

chmod u-w web/sites/default/settings.php

exit $?