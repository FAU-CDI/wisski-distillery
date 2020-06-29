#!/bin/bash
set -e

# read the lib/shared.sh
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"

# update_stack fully updates a docker-compose stack in the given location. 
function update_stack() {
    cd "$1"
    docker-compose pull
    docker-compose build --pull
    docker-compose up -d
}

log_info "=> Rebuilding and restarting 'web' stack"
update_stack "$DEPLOY_WEB_DIR"

# build and start the triplestore
log_info "=> Rebuilding and restarting 'triplestore' stack"
update_stack "$DEPLOY_TRIPLESTORE_DIR"

# build and start the triplestore
log_info "=> Rebuilding and restarting 'sql' stack"
update_stack "$DEPLOY_SQL_DIR"

# TODO: Iterate over all the instance
# and a  pull_and_update

log_info "=> System up-to-date. "