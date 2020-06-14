#!/bin/bash

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


# A global flag 'USE_DRUPAL_9' can be set to enable drupal 9 support. 
# We print out the value of the flag here
if [ -z "${USE_DRUPAL_9}" ]; then
    log_info " => Will install stable Drupal 8 (Use 'USE_DRUPAL_9=1' for Drupal 9). "
else
    log_info " => Will install experimental Drupal 9 version ('USE_DRUPAL_9' was set)"
fi

log_info " => Validing configuration"

# If the base directory already exists, we might have accidentally picked a name that already exists. 
# In that case we bail out for safety reasons. 
if [ -d "$BASE_DIR" ]; then
    echo "'$BASE_DIR' already exists. "
    echo "Aborting provisioning, please make sure you picked a unique name. "
    exit 1
fi

# Check that the apache2 config is correct. 
# This is a sanity test so that we don't randomly fail later because of bad config. 
log_info " => Checking apache configuration"
apache2ctl configtest > /dev/null

# Create a system user and group
log_info " => Creating system user and group '$SYSTEM_USER'"
addgroup --system "$SYSTEM_USER"
adduser --home "$BASE_DIR" --system --disabled-password --disabled-login --ingroup "$SYSTEM_USER" "$SYSTEM_USER"

# Make directory for the composer project to live in
log_info " => Making composer directory '$COMPOSER_DIR'"
sudo -u "$SYSTEM_USER" mkdir -p "$COMPOSER_DIR"
cd "$COMPOSER_DIR"

# Write out a new apache configuration file into /etc/apache2/sites-available. 
# We will need to substitute in some configuration directories. 
log_info " => Writing new apache configuration file"
load_template "wisski-site.conf" \
    "PUBLIC_PORT" "${PUBLIC_PORT}" \
    "WEB_DIR" "${WEB_DIR}" \
    "INSTANCE_DOMAIN" "${INSTANCE_DOMAIN}" \
    "SYSTEM_USER" "${SYSTEM_USER}" \
    "WISSKI_COMMON_PATH" "${WISSKI_COMMON_PATH}" \
    > "${APACHE_CONFIG_SITE_AVAILABLE}"

# Create a new composer project. 
log_info " => Creating composer project"
if [ -z "${USE_DRUPAL_9}" ]; then
    composer create-project 'drupal/recommended-project:^8.9.0' .
else
    composer create-project 'drupal/recommended-project:^9.0.0' .
fi
composer require drush/drush

# Randomly generate the database name and user we will configure. 
# Use the 'randompw' alias for this. 
log_info " => Generating new MySQL password"
MYSQL_PASSWORD="$(randompw)"

# Initialize the SQL database with those credentials. 
log_info " => Intializing new SQL database '${MYSQL_DATABASE}' and user '$MYSQL_USER'. "
mysql -e "CREATE DATABASE \`${MYSQL_DATABASE}\`;"
mysql -e "CREATE USER \`${MYSQL_USER}\`@localhost IDENTIFIED BY '${MYSQL_PASSWORD}';"
mysql -e "GRANT ALL PRIVILEGES ON \`${MYSQL_DATABASE}\`.* TO \`${MYSQL_USER}\`@localhost;"
mysql -e "FLUSH PRIVILEGES;"

# Generate some more random credentials, this time for drupal. 
# We again make use of the randompw alias. 
log_info " => Generating new drupal credentials"
DRUPAL_USER="admin"
DRUPAL_PASS="$(randompw)"

# Use 'drush' to run the site-installation. 
# Here we need to use the username, password and database creds we made above. 
log_info " => Running drupal installation scripts"
drush site-install standard --yes --site-name=${INSTANCE_DOMAIN} \
    --account-name=$DRUPAL_USER --account-pass=$DRUPAL_PASS \
    --db-url=mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@localhost/${MYSQL_DATABASE}

drupal_sites_permission_workaround

# Create a new repository for GraphDB. 
# Use the template for this.
log_info " => Generating new GraphDB repository '$GRAPHDB_REPO'"
load_template "graphdb-repo.ttl" "GRAPHDB_REPO" "${GRAPHDB_REPO}" "INSTANCE_DOMAIN" "${INSTANCE_DOMAIN}" | \
curl -X POST \
    http://127.0.0.1:7200/rest/repositories \
    --header 'Content-Type: multipart/form-data' \
    -F "config=@-"

# Generate a random password for the GraphDB user
log_info " => Generating a new GraphDB password"
GRAPHDB_PASSWORD="$(randompw)"

# Create the user and grant them access to the creatd database. 
log_info " => Creating GraphDB user '$GRAPHDB_USER'"
load_template "graphdb-user.json" "GRAPHDB_USER" "${GRAPHDB_USER}" "GRAPHDB_REPO" "${GRAPHDB_REPO}" | \
curl -X POST \
    "http://127.0.0.1:7200/rest/security/user/${GRAPHDB_USER}" \
    --header 'Content-Type: application/json' \
    --header 'Accept: text/plain' \
    --header "X-GraphDB-Password: $GRAPHDB_PASSWORD" \
    -d @-

# create a directory for ontologies. 
log_info " => Creating '$ONTOLOGY_DIR'"
mkdir -p "$ONTOLOGY_DIR"

# Install the Wisski packages. 
log_info " => Installing Wisski packages"
cd "$COMPOSER_DIR"

drupal_sites_permission_workaround

# install the development version when requested
if [ -z "${USE_DRUPAL_9}" ]; then
    composer require 'drupal/wisski'
else
    composer require 'drupal/wisski:2.x-dev'
fi

drupal_sites_permission_workaround
composer require drupal/inline_entity_form

drupal_sites_permission_workaround
composer require drupal/imagemagick

drupal_sites_permission_workaround
composer require drupal/image_effects

drupal_sites_permission_workaround
composer require drupal/colorbox

log_info " => Installation is now technically complete. "
log_ok "Some things below may fail. If that is the case, run: "
log_ok "$ a2ensite \"${INSTANCE_DOMAIN}\""
log_ok "$ systemctl reload apache2"
log_ok "$ $SCRIPT_DIR/shell.sh $SLUG"
log_ok "Your installation details are as follows:"
function printdetails() {
    echo "URL:                  http://$INSTANCE_DOMAIN"
    echo "Username:             $DRUPAL_USER"
    echo "Password:             $DRUPAL_PASS"
    log_info " => Your GraphDB details (for WissKI Salz) are: "
    echo "Read URL:             http://127.0.0.1:7200/repositories/$GRAPHDB_REPO"
    echo "Write URL:            http://127.0.0.1:7200/repositories/$GRAPHDB_REPO/statements"
    echo "Writable:             yes"
    echo "Default Graph URI:    http://$INSTANCE_DOMAIN/#"
    echo "Ontology Paths:       (empty)"
    echo "SameAs property:      http://www.w3.org/2002/07/owl#sameAs"
}
printdetails

function alldetails() {
    echo "# Automatically generated WissKi details"
    echo "# generated $(date -u +"%Y-%m-%dT%H:%M:%SZ")"
    echo "SLUG=$SLUG"
    echo "INSTANCE_DOMAIN=$INSTANCE_DOMAIN"
    echo "# System"
    echo "SYSTEM_USER=$SYSTEM_USER"
    echo "# Drupal"
    echo "DRUPAL_USER=$DRUPAL_USER"
    echo "DRUPAL_PASSWORD=$DRUPAL_PASSWORD"
    echo "# MySQL"
    echo "MYSQL_USER=$MYSQL_USER"
    echo "MYSQL_PASSWORD=$MYSQL_PASSWORD"
    echo "MYSQL_DATABASE=$MYSQL_DATABASE"
    echo "# GraphDB"
    echo "GRAPHDB_USER=$GRAPHDB_USER"
    echo "GRAPHDB_PASSWORD=$GRAPHDB_PASSWORD"
    echo "GRAPHDB_REPO=$GRAPHDB_REPO"
}

# put installation details in ENV_FILE
log_info " => Storing installation details in $ENV_FILE"
alldetails > "$ENV_FILE"
chown "$SYSTEM_USER:$SYSTEM_USER" "$ENV_FILE"
chmod o-r "$ENV_FILE"

# Enable the WissKI modules. 
log_info " => Enable Wisski modules"
drush pm-enable --yes wisski_core wisski_linkblock wisski_pathbuilder wisski_adapter_sparql11_pb wisski_salz
drupal_sites_permission_workaround

# Because of a regresssion in EasyRDF and Tomcat, we need to manually patch EasyRDF
EASYRDF_RESPONSE="$COMPOSER_DIR/vendor/easyrdf/easyrdf/lib/EasyRdf/Http/Response.php"
log_info " => Patching '$EASYRDF_RESPONSE'"
load_template "easyrdf.patch" | patch "$EASYRDF_RESPONSE"

# Finally enable the apache2 config. 
# And then reload to start serving it. 
log_info " => Enabling and reloading apache configuration"
a2ensite "${INSTANCE_DOMAIN}"
systemctl reload apache2

# and done!
log_info " => Finished, your Drupal details are: "
printdetails
