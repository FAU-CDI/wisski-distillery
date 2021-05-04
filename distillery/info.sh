#!/bin/bash
set -e

# read the lib/shared.sh and read the slug argument. 
DISABLE_LOG=1
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"
DISABLE_LOG=0
require_slug_argument


# if the site doesn't exist, I can't open a shell. 
if ! sql_bookkeep_exists "$SLUG"; then
    log_error "=> Site '$SLUG' does not exist in bookeeping table. "
    echo "I can't show info about it. "
    exit 1
fi;

# Read everything from the database
read -r INSTANCE_BASE_DIR MYSQL_DATABASE MYSQL_USER GRAPHDB_REPO GRAPHDB_USER GRAPHDB_PASSWORD <<< "$(sql_bookkeep_load "${SLUG}" "filesystem_base,sql_database,sql_user,graphdb_repository,graphdb_user,graphdb_password" | tail -n +2)"

GRAPHDB_HEADER="$(printf "%s:%s" "$GRAPHDB_USER" "$GRAPHDB_PASSWORD" | base64 -w 0)"

# read sql configuration
cd "$INSTANCE_BASE_DIR"
read -r  SQL_DATABASE SQL_USER SQL_PASS SQL_OTHER <<< "$(docker-compose exec barrel drush sql:conf --format=tsv --show-passwords)"

echo "=================================================================================="
echo "URL:                  http://$INSTANCE_DOMAIN"
echo "Base directory:       ${INSTANCE_BASE_DIR}"
log_info " => Your GraphDB details (for WissKI Salz) are: "
echo "Read URL:             http://triplestore:7200/repositories/$GRAPHDB_REPO"
echo "Write URL:            http://triplestore:7200/repositories/$GRAPHDB_REPO/statements"
echo "Username:             $GRAPHDB_USER"
echo "Password:             $GRAPHDB_PASSWORD"
echo "Authorization Header: $GRAPHDB_HEADER"
echo "Writable:             yes"
echo "Default Graph URI:    http://$INSTANCE_DOMAIN/#"
echo "Ontology Paths:       (empty)"
echo "SameAs property:      http://www.w3.org/2002/07/owl#sameAs"
log_info " => Your SQL detsils are: "
echo "SQL Database:         $SQL_DATABASE"
echo "SQL Username:         $SQL_USER"
echo "SQL Password:         $SQL_PASS"
