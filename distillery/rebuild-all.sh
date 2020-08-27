#!/bin/bash
set -e

# read the lib/shared.sh and read the slug argument. 
DISABLE_LOG=1
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"
DISABLE_LOG=0

# update all the instances
for slug in $(sql_bookkeep_list); do
    log_info "=> /bin/bash $DIR/rebuild.sh '$slug'"
    /bin/bash "$DIR/rebuild.sh" "$slug";
done

