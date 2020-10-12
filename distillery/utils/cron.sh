#!/bin/bash

# This utility script can be used to run all cron tasks. 

cd /var/www/data/project || exit 1
export PATH=/var/www/data/project/vendor/bin:$PATH

drush core-cron