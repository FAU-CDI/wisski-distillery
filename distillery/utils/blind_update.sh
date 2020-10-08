#!/bin/bash

# This utility script can be used to blindly update all dependencies to their latest versions. 
# It does not perform any checking whatsoever. 

cd /var/www/data/project || exit 1

# composer install updates
chmod u+rw web/sites/default/
composer update

# update the dabatabase
drush -y updatedb