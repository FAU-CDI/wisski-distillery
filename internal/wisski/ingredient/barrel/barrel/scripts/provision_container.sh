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

DRUPAL_VERSION="$1"
echo "DRUPAL_VERSION=$DRUPAL_VERSION"
shift 1

WISSKI_VERSION="$1"
echo "WISSKI_VERSION=$WISSKI_VERSION"
shift 1

log_info " => Preparing installation environment"
BASE_DIR="/var/www/data"
COMPOSER_DIR="$BASE_DIR/project"
WEB_DIR="$COMPOSER_DIR/web"
ONTOLOGY_DIR="$WEB_DIR/sites/default/files/ontology"

log_info " => Creating '$COMPOSER_DIR'"
mkdir -p "$COMPOSER_DIR"
cd "$COMPOSER_DIR"

# workaround for making the drupal sites directory writable
function drupal_sites_permission_workaround() {
    chmod -R u+w "$WEB_DIR/sites/" || true
}

# install a module with composer and enable it with drush
# Example:
#
# composer_install_and_enable << EOF
# drupal/some_module:1.23 some_module
# drupal/other_module:2.34
# EOF
# 
# Will install both modules, but only enable the first one.
function composer_install_and_enable() {
    while IFS= read -r line; do
        echo "$line" | (
            read composer drush;
            drupal_sites_permission_workaround
            composer require "$composer"
            if [ -n "$drush" ]; then
                drush pm-enable --yes "$drush"
            fi
        )
    done
}

function try_variants() {
    for var in "$@"
    do
        if composer require --dry-run "$var" > /dev/null 2>&1; then
            composer require "$var"
            return 0;
        fi
    done

    return 1;
}


# Create a new composer project. 
log_info " => Creating composer project"
if [ -z "${DRUPAL_VERSION}" ]; then
    composer --no-interaction create-project 'drupal/recommended-project:^9.0.0' .
else
    composer --no-interaction create-project "drupal/recommended-project:$DRUPAL_VERSION" .
fi

# needed for composer > 2.2
composer --no-interaction config allow-plugins true

# Install drush so that we can automate a lot of things
log_info " => Installing 'drush'"
try_variants 'drush/drush' 'drush/drush:^12' 'drush/drush:^11' || (echo "No version of Drush is installable" && false)

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

# Install some additional modules
# These neeed to go before WissKI because some are WissKI dependencies

log_info " => Installing and enabling modules"
composer_install_and_enable << EOF
drupal/inline_entity_form:^1.0@RC
drupal/imagemagick
drupal/image_effects
drupal/colorbox
drupal/devel:^4.1 devel
drupal/geofield:^1.40 geofield
drupal/geofield_map:^2.85 geofield_map
drupal/imce:^2.4 imce
drupal/remove_generator:^2.0 remove_generator
EOF


# Install the Wisski packages. 
log_info " => Installing Wisski packages"
cd "$COMPOSER_DIR"

# install the development version when requested
if [ -z "${WISSKI_VERSION}" ]; then
    composer require 'drupal/wisski'
else
    composer require "drupal/wisski:$WISSKI_VERSION"
fi

# Install dependencies of WissKI
log_info " => Installing and patching Wisski dependencies"
pushd "$WEB_DIR/modules/contrib/wisski"
composer install

# Patch EasyRDF (for now)
EASYRDF_RESPONSE="./vendor/easyrdf/easyrdf/lib/EasyRdf/Http/Response.php"
if [ -f "$EASYRDF_RESPONSE" ]; then
    patch -N "$EASYRDF_RESPONSE" < "/patch/easyrdf.patch"
fi
popd



log_info " => Enable Wisski modules"
drush pm-enable --yes wisski_core wisski_linkblock wisski_pathbuilder wisski_adapter_sparql11_pb wisski_salz
drupal_sites_permission_workaround

log_info " => Setting up WissKI Salz Adapter"
drush php:script /wisskiutils/create_adapter.php "$INSTANCE_DOMAIN" "$GRAPHDB_REPO" "$GRAPHDB_HEADER"

log_info " => Updating TRUSTED_HOST_PATTERNS in settings.php"

/bin/bash /wisskiutils/set_trusted_host.sh

log_info " => Running initial cron"
drush core-cron

log_info " => Provisioning is now complete. "
log_ok "Your installation details are as follows:"
function printdetails() {
    echo "URL:                  http://$INSTANCE_DOMAIN"
    echo "Username:             $DRUPAL_USER"
    echo "Password:             $DRUPAL_PASS"
}
printdetails

exit 0