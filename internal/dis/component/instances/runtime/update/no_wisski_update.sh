#!/bin/bash
set -e

# Performs an update of this instance, leaving the WissKI package untouched.
# It is recommended to run update-composer-json.php prior to this script.

cd "/var/www/data/project"

# Fix permissions
chmod -R 755 web/modules web/themes web/profiles web/core
chmod -R 755 web/sites/default/files
chmod 755 web/sites/default
chmod 644 web/sites/*/settings*.php
chmod 444 web/.htaccess

# Run the actual update
php /runtime/update/composer-no-wisski-update.php --with-all-dependencies

# Clear caches and apply database updates
drush cr 
drush updatedb --yes
drush cr