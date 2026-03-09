#!/bin/bash

# chown the volumes to make sure they can be read and written by the limited user
chown www-data:www-data /var/www
chown www-data:www-data /var/www/.composer
chown www-data:www-data /var/www/data/

# start up dropbear
/ssh/start.sh &

# run the original entrypoint
docker-php-entrypoint "$@"