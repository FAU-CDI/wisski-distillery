#!/bin/bash
set -e

# This is a library file. 
# It should be 'source'd only, if it is not we bail out here. 
if [[ "$0" = "$BASH_SOURCE" ]]; then
   echo "This file should not be executed directly, it should be 'source'd only. "
   exit 1;
fi

# Set a few variables to point to the debian frontend
export DEBIAN_FRONTEND=noninteractive

# This file just sets a few utility functions to be used by the code. 
# randompw generates a random password as per the configuration file. 
alias randompw="cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w $PASSWORD_LENGTH | head -n 1"

# update_stack fully updates a docker-compose stack in the given location. 
function update_stack() {
    cd "$1"
    docker-compose pull
    docker-compose build --pull
    docker-compose up -d
}