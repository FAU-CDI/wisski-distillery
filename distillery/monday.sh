#!/bin/bash
set -e

# read the lib/shared.sh and read the slug argument. 
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"

# Read the 'GRAPHDB_ZIP' argument from the command line. 
# If it's not set, throw an error. 
GRAPHDB_ZIP=$1
if [ -z "$GRAPHDB_ZIP" ]; then
    log_error "Usage: monday.sh GRAPHDB_ZIP"
    exit 1;
fi;


# Backup
log_info " => Running backup, this will take a long time"
bash backup.sh

# system install
log_info " => Reinstalling system"
bash system_install.sh "$GRAPHDB_ZIP"

# rebuild all the systems
log_info " => Rebuilding all instances"
bash rebuild-all.sh

# perform all the blind updates
log_info " => Performing updates"
bash blind-update-all.sh

log_info " => Done, have a great week"
