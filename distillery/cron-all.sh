#!/bin/bash
set -e

# read the lib/shared.sh and read the slug argument. 
DISABLE_LOG=1
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"
unset DISABLE_LOG

# update all the instances
for slug in $(sql_bookkeep_list); do
    read -r INSTANCE_BASE_DIR <<< "$(sql_bookkeep_load "${slug}" "filesystem_base" | tail -n +2)"
    log_info "=> Runnning cron for '$slug'"
    cd "$INSTANCE_BASE_DIR"
    docker-compose exec barrel /bin/bash /utils/cron.sh
done

