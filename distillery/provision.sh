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

# check if the site exists
if sql_bookkeep_exists "$SLUG"; then
    log_error "=> Site '$SLUG' already exists in bookeeping table. "
    echo "Refusing to work"
    exit 1;
fi

# Randomly generate the database name and user we will configure. 
# Use the 'randompw' alias for this. 
log_info " => Generating new MySQL password"
MYSQL_PASSWORD="$(randompw)"

# Initialize the SQL database with those credentials. 
log_info " => Intializing new SQL database '${MYSQL_DATABASE}' and user '$MYSQL_USER'. "
dockerized_mysql -e "CREATE DATABASE \`${MYSQL_DATABASE}\`;"
dockerized_mysql -e "CREATE USER \`${MYSQL_USER}\`@'%' IDENTIFIED BY '${MYSQL_PASSWORD}';"
dockerized_mysql -e "GRANT ALL PRIVILEGES ON \`${MYSQL_DATABASE}\`.* TO \`${MYSQL_USER}\`@\`%\`;"
dockerized_mysql -e "FLUSH PRIVILEGES;"

# Create a new repository for GraphDB. 
# Use the template for this.
log_info " => Generating new GraphDB repository '$GRAPHDB_REPO'"
load_template "repository/graphdb-repo.ttl" "GRAPHDB_REPO" "${GRAPHDB_REPO}" "INSTANCE_DOMAIN" "${INSTANCE_DOMAIN}" | \
curl -X POST \
    http://127.0.0.1:7200/rest/repositories \
    --header 'Content-Type: multipart/form-data' \
    -F "config=@-"

# Generate a random password for the GraphDB user
log_info " => Generating a new GraphDB password"
GRAPHDB_PASSWORD="$(randompw)"

# Create the user and grant them access to the creatd database. 
log_info " => Creating GraphDB user '$GRAPHDB_USER'"
load_template "repository/graphdb-user.json" "GRAPHDB_USER" "${GRAPHDB_USER}" "GRAPHDB_REPO" "${GRAPHDB_REPO}" | \
curl -X POST \
    "http://127.0.0.1:7200/rest/security/user/${GRAPHDB_USER}" \
    --header 'Content-Type: application/json' \
    --header 'Accept: text/plain' \
    --header "X-GraphDB-Password: $GRAPHDB_PASSWORD" \
    -d @-

log_info " => Creating local directory '$INSTANCE_BASE_DIR'"
mkdir -p "$INSTANCE_BASE_DIR"
mkdir -p "$INSTANCE_DATA_DIR"
mkdir -p "$INSTANCE_DATA_DIR/.composer"
mkdir -p "$INSTANCE_DATA_DIR/data"

# Generate some more random credentials, this time for drupal. 
# We again make use of the randompw alias. 
log_info " => Generating new drupal credentials"
DRUPAL_USER="admin"
DRUPAL_PASS="$(randompw)"

# TODO: copy over docker-compose into the right directory
log_info " => Creating instance directory"
install_resource_dir "compose/runtime" "$INSTANCE_BASE_DIR"

# Log all the details into the bookeeping database
log_info "=> Storing configuration in bookkeeping table"
sql_bookkeep_insert \
    "slug,filesystem_base,sql_database,sql_user,sql_password,graphdb_repository,graphdb_user,graphdb_password" \
    "\"${SLUG}\",\"${INSTANCE_BASE_DIR}\",\"${MYSQL_DATABASE}\",\"${MYSQL_USER}\",\"${MYSQL_PASSWORD}\",\"${GRAPHDB_REPO}\",\"${GRAPHDB_USER}\",\"${GRAPHDB_PASSWORD}\""

log_info " => Writing configuration file"
load_template "runtime-config/environment" \
    "REAL_PATH" "${INSTANCE_DATA_DIR}" \
    "VIRTUAL_HOST" "${INSTANCE_DOMAIN}" \
    "LETSENCRYPT_HOST" "${LETSENCRYPT_HOST}" \
    "LETSENCRYPT_EMAIL" "${LETSENCRYPT_EMAIL}" \
    > "$INSTANCE_BASE_DIR/.env"


log_info " => Running and building image"
cd "$INSTANCE_BASE_DIR"
docker-compose build --pull
docker-compose pull

log_info " => Running provision script"
docker-compose run --rm runtime /bin/bash -c "sudo PATH=\$PATH -u www-data /bin/bash /provision_container.sh \
    \"${INSTANCE_DOMAIN}\" \
    \"${MYSQL_DATABASE}\" \"${MYSQL_USER}\" \"${MYSQL_PASSWORD}\" \
    \"${GRAPHDB_REPO}\" \"${GRAPHDB_USER}\" \"${GRAPHDB_PASSWORD}\" \
    \"${DRUPAL_USER}\" \"${DRUPAL_PASS}\" \
    \"${USE_DRUPAL_9}\""


log_info " => Starting container"
docker-compose up -d

