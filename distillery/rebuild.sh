#!/bin/bash
set -e

# read the lib/shared.sh and read the slug argument. 
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"
require_slug_argument


# if the site doesn't exist, I can't open a shell. 
if ! sql_bookkeep_exists "$SLUG"; then
    log_error "=> Site '$SLUG' does not exist in bookeeping table. "
    echo "I can't rebuild it. "
    exit 1
fi;

# Read everything from the database
read -r INSTANCE_BASE_DIR MYSQL_DATABASE MYSQL_USER GRAPHDB_REPO GRAPHDB_USER <<< "$(sql_bookkeep_load "${SLUG}" "filesystem_base,sql_database,sql_user,graphdb_repository,graphdb_user" | tail -n +2)"

log_info " => Touching authorized_keys file"
touch "$INSTANCE_BASE_DIR/data/authorized_keys"

log_info " => Updating compose files"
install_resource_dir "compose/barrel" "$INSTANCE_BASE_DIR"

log_info "=> Rebuilding and restarting '$INSTANCE_BASE_DIR'"
update_stack "$INSTANCE_BASE_DIR"