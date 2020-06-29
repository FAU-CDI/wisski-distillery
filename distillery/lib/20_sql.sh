#!/bin/bash
set -e

# This is a library file. 
# It should be 'source'd only, if it is not we bail out here. 
if [[ "$0" = "$BASH_SOURCE" ]]; then
   echo "This file should not be executed directly, it should be 'source'd only. "
   exit 1;
fi

###
### General SQL functions
###

# wait_for_sql waits for the sql database to be ready
function wait_for_sql() {
    log_info "Waiting for database to start ..."
    _wait_for_sql_internal
    log_ok "Database responded to query "
}

function _wait_for_sql_internal() {
    timeout=30
    times=1
    until dockerized_mysql -e 'show databases;' > /dev/null 2>&1; do
        ((times=times+1))
        if [ "$times" -gt $timeout ]; then
            log_error "Database timed out after ${timeout} seconds(s). "
        fi;
        echo -n "."
        sleep 1
    done
}

# 'dockerized_mysql' runs an sql command in the sql docker container
function dockerized_mysql() {
    pushd "$DEPLOY_SQL_DIR" > /dev/null
    docker exec -i `docker-compose ps -q sql` mysql "$@"
    retval=$?
    popd > /dev/null
    return $retval
}

# 'dockerized_mysql' runs an sql command in the sql docker container interactively
function dockerized_mysql_interactive() {
    pushd "$DEPLOY_SQL_DIR" > /dev/null
    docker exec -ti `docker-compose ps -q sql` mysql "$@"
    retval=$?
    popd > /dev/null
    return $retval
}

###
### Bookkeeping sql
###


# 'sql_bookkeep_exists' checks if a given site already exists
function sql_bookkeep_exists() {
    COUNT=$(dockerized_mysql -D "$DISTILLERY_BOOKKEEPING_DATABASE" -e "SELECT COUNT(*) FROM \`$DISTILLERY_BOOKKEEPING_TABLE\` WHERE slug=\"$1\"; "  | tail -n +2)
    if [ "$COUNT" = "1" ]; then
        return 0;
    else
        return 1;
    fi;
}

# 'sql_bookkeep_insert' inserts a new pair of values into the sql database
function sql_bookkeep_insert() {
    dockerized_mysql -D "$DISTILLERY_BOOKKEEPING_DATABASE" -e "INSERT INTO \`$DISTILLERY_BOOKKEEPING_TABLE\`($1) VALUES ($2) ;"
}

# 'sql_bookkeep_delete' removes a slug into the bookeeping table
function sql_bookeep_delete() {
    dockerized_mysql -D "$DISTILLERY_BOOKKEEPING_DATABASE" -e "DELETE FROM \`$DISTILLERY_BOOKKEEPING_TABLE\` WHERE slug=\"$1\";"
}

# 'sql_bookkeep_load' reads from the bookkeeping table
function sql_bookkeep_load() {
    dockerized_mysql -D "$DISTILLERY_BOOKKEEPING_DATABASE" -e "SELECT $2 FROM \`$DISTILLERY_BOOKKEEPING_TABLE\` WHERE slug=\"$1\";"
}

# 'sql_bookkeep_list' lists all slugs from the bookkeeping table
function sql_bookkeep_list() {
    dockerized_mysql -D "$DISTILLERY_BOOKKEEPING_DATABASE" -e "SELECT slug FROM \`$DISTILLERY_BOOKKEEPING_TABLE\`; " | tail -n +2
}