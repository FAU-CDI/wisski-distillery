#!/bin/bash

# This utility script can be used to blindly update all dependencies to their latest versions. 
# It does not perform any checking whatsoever. 

# update the main modules
cd /var/www/data/project || exit 1
chmod u+rw web/sites/default/
composer update

# update the db
drush -y updatedb

# update the wisski dependencies
cd /var/www/data/project/web/modules/contrib/wisski || exit 1
composer update