#!/bin/bash

# settings_php_get.sh name
# Gets the 'settings_php_get.php' setting 'name' as json-encoded value, or null when it does not exist.

NAME=$1

if [ -z "$NAME" ]; then
    echo "Usage: get_settings_setting.sh NAME"
    exit 1
fi;

echo "$NAME" | drush php:eval '
  use \Drupal\Core\Site\Settings;
  $name=trim(file_get_contents("php://stdin"));
  echo json_encode(Settings::get($name));
';