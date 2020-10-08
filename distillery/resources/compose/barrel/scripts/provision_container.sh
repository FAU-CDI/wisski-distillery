#!/bin/bash
set -e

function log_info() {
   echo -e "\033[1m$1\033[0m"
}

function log_ok() {
   echo -e "\033[0;32m$1\033[0m"
}

log_info " => Reading configuration variables"

INSTANCE_DOMAIN="$1"
echo "INSTANCE_DOMAIN=$INSTANCE_DOMAIN"
shift 1

MYSQL_DATABASE="$1"
echo "MYSQL_DATABASE=$MYSQL_DATABASE"
MYSQL_USER="$2"
echo "MYSQL_USER=$MYSQL_USER"
MYSQL_PASSWORD="$3"
echo "MYSQL_PASSWORD=$MYSQL_PASSWORD"

shift 3

GRAPHDB_REPO="$1"
echo "GRAPHDB_REPO=$GRAPHDB_REPO"
GRAPHDB_USER="$2"
echo "GRAPHDB_USER=$GRAPHDB_USER"
GRAPHDB_PASSWORD="$3"
echo "GRAPHDB_PASSWORD=$GRAPHDB_PASSWORD"
shift 3

GRAPHDB_HEADER="$(printf "%s:%s" "$GRAPHDB_USER" "$GRAPHDB_PASSWORD" | base64 -w 0)"

DRUPAL_USER="$1"
echo "DRUPAL_USER=$DRUPAL_USER"
DRUPAL_PASS="$2"
echo "DRUPAL_PASS=$DRUPAL_PASS"
shift 2

USE_DRUPAL_9="$1"
echo "USE_DRUPAL_9=$USE_DRUPAL_9"
shift 1

log_info " => Preparing installation environment"
BASE_DIR="/var/www/data"
COMPOSER_DIR="$BASE_DIR/project"
WEB_DIR="$COMPOSER_DIR/web"
ONTOLOGY_DIR="$WEB_DIR/sites/default/files/ontology"

log_info " => Creating '$COMPOSER_DIR'"
mkdir -p "$COMPOSER_DIR"
cd "$COMPOSER_DIR"

function drupal_sites_permission_workaround() {
    chmod -R u+w "$WEB_DIR/sites/" || true
}

# Create a new composer project. 
log_info " => Creating composer project"
if [ -z "${USE_DRUPAL_9}" ]; then
    composer create-project 'drupal/recommended-project:^8.9.0' .
else
    composer create-project 'drupal/recommended-project:^9.0.0' .
fi

# Install drush so that we can automate a lot of things
log_info " => Installing 'drush'"
composer require drush/drush

# Use 'drush' to run the site-installation. 
# Here we need to use the username, password and database creds we made above. 
log_info " => Running drupal installation scripts"
drush site-install standard --yes --site-name=${INSTANCE_DOMAIN} \
    --account-name=$DRUPAL_USER --account-pass=$DRUPAL_PASS \
    --db-url=mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@sql/${MYSQL_DATABASE}
drupal_sites_permission_workaround

# create a directory for ontologies. 
log_info " => Creating '$ONTOLOGY_DIR'"
mkdir -p "$ONTOLOGY_DIR"

# Install the Wisski packages. 
log_info " => Installing Wisski packages"
cd "$COMPOSER_DIR"

# install the development version when requested
if [ -z "${USE_DRUPAL_9}" ]; then
    composer require 'drupal/wisski'
else
    composer require 'drupal/wisski:2.x-dev'
fi

# Install dependencies of WissKI
log_info " => Installing and patching Wisski dependencies"
pushd "$WEB_DIR/modules/contrib/wisski"
composer install

# Patch EasyRDF (for now)
EASYRDF_RESPONSE="./vendor/easyrdf/easyrdf/lib/EasyRdf/Http/Response.php"
patch -N "$EASYRDF_RESPONSE" < "/patch/easyrdf.patch"
popd

drupal_sites_permission_workaround
composer require drupal/inline_entity_form

drupal_sites_permission_workaround
composer require drupal/imagemagick

drupal_sites_permission_workaround
composer require drupal/image_effects

drupal_sites_permission_workaround
composer require drupal/colorbox

log_info " => Enable Wisski modules"
drush pm-enable --yes wisski_core wisski_linkblock wisski_pathbuilder wisski_adapter_sparql11_pb wisski_salz
drupal_sites_permission_workaround

log_info " => Provisioning is now complete. "
log_ok "Your installation details are as follows:"
function printdetails() {
    echo "URL:                  http://$INSTANCE_DOMAIN"
    echo "Username:             $DRUPAL_USER"
    echo "Password:             $DRUPAL_PASS"
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
}
printdetails