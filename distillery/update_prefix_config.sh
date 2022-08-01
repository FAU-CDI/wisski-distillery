#!/bin/bash
set -e

# read the lib/shared.sh
DISABLE_LOG=1
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"
DISABLE_LOG=0

log_info " => Writing prefix configuration"

echo -n "# Prefix configuration, last updated on" | tee "$DEPLOY_PREFIX_CONFIG"
date | tee -a "$DEPLOY_PREFIX_CONFIG"

# update all the instances
for slug in $(sql_bookkeep_list); do
    INSTANCE_DOMAIN="http://$(compute_instance_domain "$slug")"
    echo "$INSTANCE_DOMAIN:" | tee -a "$DEPLOY_PREFIX_CONFIG"
    
    read -r INSTANCE_BASE_DIR MYSQL_DATABASE MYSQL_USER GRAPHDB_REPO GRAPHDB_USER GRAPHDB_PASSWORD <<< "$(sql_bookkeep_load "${slug}" "filesystem_base,sql_database,sql_user,graphdb_repository,graphdb_user,graphdb_password" | tail -n +2)"

    pushd "$INSTANCE_BASE_DIR" > /dev/null

    INSTANCE_PREFIX_FILE="$(compute_instance_prefixfile "$INSTANCE_BASE_DIR" )"
    if [ -f "$INSTANCE_PREFIX_FILE" ]; then
        cat "$INSTANCE_PREFIX_FILE" | tee -a "$DEPLOY_PREFIX_CONFIG"
    fi
    
    docker-compose exec barrel /user_shell.sh -c "drush php:script /wisskiutils/list_uri_prefixes.php" | tee -a "$DEPLOY_PREFIX_CONFIG"
    popd > /dev/null
done

log_info " => Restarting resolver"

cd "$DEPLOY_RESOLVER_DIR"
docker-compose restart