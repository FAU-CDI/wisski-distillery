#!/bin/bash

chmod 777 web/sites/default
chmod 666 web/sites/default/*settings.php
chmod 666 web/sites/default/*services.yml

composer require 'drupal/core-recommended:^9' 'drupal/core-composer-scaffold:^9' 'drupal/core-project-message:^9' --update-with-dependencies --no-update
composer update
composer require 'drupal/wisski'

drush updatedb

chmod 755 web/sites/default
chmod 644 web/sites/default/*settings.php
chmod 644 web/sites/default/*services.yml