# To install a new system:

# This script will provision a new Drupal instance and make it available to apache. 
# Usage: sudo ./provision.sh $SLUG
# In case the installation fails, it will bail out and leave you with an incomplete installation. 
# To delete an incomplete installation, use the ./purge.sh script, or try fixing the error manually. 
set -e

# read the lib/shared.sh and read the slug argument. 
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"
require_slug_argument

# wait for sql to be awake
wait_for_sql

while true; do
    log_info " => I'm about to delete the '$SLUG' site from this system. "
    read -p "This can not be undone. Please type 'y' to continue: " yn
    case $yn in
        [Yy]* ) break;;
        * ) echo "Abort. "; exit 1;;
    esac
done

# check if the site exists
if ! sql_bookkeep_exists "$SLUG"; then
    log_error "=> Site '$SLUG' does not exist in bookeeping table. "
    echo "I'll try to cleanup with the current defaults. "
    echo "This may or may not work. "
else
    # Read all the configuration from the database
    log_info " => Reading components from database"
    read -r INSTANCE_BASE_DIR MYSQL_DATABASE MYSQL_USER GRAPHDB_REPO GRAPHDB_USER <<< "$(sql_bookkeep_load "${SLUG}" "filesystem_base,sql_database,sql_user,graphdb_repository,graphdb_user" | tail -n +2)"
fi

# stop the running system container
if [ -d "$INSTANCE_BASE_DIR" ] ; then
    log_info "=> Stopping running system"
    cd "$INSTANCE_BASE_DIR"
    docker-compose down -v || true
fi;

cd

# delete the mysql database. 
log_info " => Deleting MySQL database '$MYSQL_DATABASE' and user '$MYSQL_USER'. "
dockerized_mysql -e "DROP DATABASE IF EXISTS \`${MYSQL_DATABASE}\`;"
dockerized_mysql -e "DROP USER IF EXISTS \`${MYSQL_USER}\`@\`%\`;"
dockerized_mysql -e "FLUSH PRIVILEGES;"

# Clear the GraphDB repository. 
log_info " => Deleting GraphDB user '$GRAPHDB_USER'"
curl -X DELETE http://127.0.0.1:7200/rest/security/user/$GRAPHDB_USER/

log_info " => Deleting GraphDB repository '$GRAPHDB_REPO'"
curl -X DELETE http://127.0.0.1:7200/rest/repositories/$GRAPHDB_REPO/

# Delete the directory
log_info " => Deleting '$INSTANCE_BASE_DIR'"
rm -rf "$INSTANCE_BASE_DIR"

log_info " => Clearing bookkeeping record"
sql_bookeep_delete "$SLUG" || true

log_info " => '$SLUG' has been purged. "
