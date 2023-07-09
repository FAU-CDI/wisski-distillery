#!/bin/bash
set -e

# temporarily extend permissions
chmod 777 web/sites/default
chmod 666 web/sites/default/*settings.php
chmod 666 web/sites/default/*services.yml

# update the core itself
composer require 'drupal/internal/core-recommended:^9' 'drupal/internal/core-composer-scaffold:^9' 'drupal/internal/core-project-message:^9' --update-with-dependencies --no-update
composer update
composer require 'drupal/wisski'

# update requirements for wisski!
pushd web/modules/contrib/wisski || exit 1
composer update
popd || exit 1

# run the update and clear the cache!
drush updatedb --yes

# and reset everything back to normal
chmod 755 web/sites/default
chmod 644 web/sites/default/*settings.php
chmod 644 web/sites/default/*services.yml