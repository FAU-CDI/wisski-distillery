#!/bin/bash

# This script is used to start a user shell inside the docker container. 
cd "/var/www/data/project"
sudo -u www-data "PATH=/var/www/data/project/vendor/bin:$PATH" /bin/bash "$@"