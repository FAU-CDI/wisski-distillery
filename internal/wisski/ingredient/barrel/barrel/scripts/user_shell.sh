#!/bin/bash
set -e

# This script is used to start a user shell inside the docker container. 
cd "/var/www/data/project"
export "PATH=/var/www/data/project/vendor/bin:$PATH"

if [ "$USER" = "www-data" ]; then
    /bin/bash "$@"
else
    sudo -u www-data  /bin/bash "$@"
fi;