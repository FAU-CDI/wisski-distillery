#!/bin/bash
set -e

# read user
VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: use_wisski.sh VERSION"
    exit 1
fi

# update the main modules
cd /var/www/data/project
chmod u+rw web/sites/default/
composer require "drupal/wisski:$VERSION"

# update the wisski dependencies
pushd /var/www/data/project/web/modules/contrib/wisski
composer update
popd

# update the db
drush -y updatedb
