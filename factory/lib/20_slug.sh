#!/bin/bash
set -e

# This is a library file. 
# It should be 'source'd only, if it is not we bail out here. 
if [[ "$0" = "$BASH_SOURCE" ]]; then
   echo "This file should not be executed directly, it should be 'source'd only. "
   exit 1;
fi

# This file reads a single slug command line option. 
# This is validated when 'require_slug_argument' is called. 

function require_slug_argument() {
    # The 'SLUG' argument must be a valid slug. 
    if ! is_valid_slug "$SLUG"; then
        log_error "Argument 'SLUG' is missing or not a valid slug. ";
        log_info "Please provide it via the command line. ";
        exit 1;
    fi

    log_info " => Deriving configuration for '$SLUG'. "
    echo "Domain Name:          $INSTANCE_DOMAIN"
    echo "Base Directory:       $BASE_DIR"
    echo "System User:          $SYSTEM_USER"
    echo "MySQL User:           $MYSQL_USER"
    echo "MySQL Database:       $MYSQL_DATABASE"
    echo "GraphDB User:         $GRAPHDB_USER"
    echo "GraphDB Repository:   $GRAPHDB_REPO"
}

# Read the slug argument. 
# We also read it in for scripts where it is not required, and will only use it if that is the case. 
SLUG="$1"

# Compute the domain name for this instance.
# Also lowercase the domain name for consistency. 
INSTANCE_DOMAIN="$SLUG.$DEFAULT_DOMAIN"
INSTANCE_DOMAIN="$(echo "$INSTANCE_DOMAIN" | tr '[:upper:]' '[:lower:]')"

# Next we need a username base. 
# This will be used as a username across the system (linux), MySQL and GraphDB. 
# For this we can only allow [0-9a-zA-Z-], hence we have to escape. 
# In most cases, the only characters that require escaping are '.'s. 
# Hence we replace '.' with '-'s.
# We replace the other two characters that require escaping (_ and -)s with --u and --s respectively. 
# Because no two dots can ever follow each other in the INSTANCE_DOMAIN, this is guaranteed collision free. 
# We also have to do the '-' replacement first, to prevent escaped other characters from being escaped twice. 
USERNAME_BASE="$SLUG"
USERNAME_BASE="${USERNAME_BASE//-/--d}"
USERNAME_BASE="${USERNAME_BASE//_/--u}"
USERNAME_BASE="${USERNAME_BASE//./-}"

# Generate the user and database names for the various systems
SYSTEM_USER="${SYSTEM_USER_PREFIX}${USERNAME_BASE}"
MYSQL_USER="${MYSQL_USER_PREFIX}${USERNAME_BASE}"
MYSQL_DATABASE="${MYSQL_DATABASE_PREFIX}${USERNAME_BASE}"
GRAPHDB_USER="${GRAPHDB_USER_PREFIX}${USERNAME_BASE}"
GRAPHDB_REPO="${GRAPHDB_REPO_PREFIX}${USERNAME_BASE}"

# Compute the base directory for the files that will live on disk. 
BASE_DIR="$DRUPAL_ROOT/$INSTANCE_DOMAIN"
ENV_FILE="$BASE_DIR/wisski-env"
COMPOSER_DIR="$BASE_DIR/project"
WEB_DIR="$COMPOSER_DIR/web"
ONTOLOGY_DIR="$WEB_DIR/sites/default/files/ontology"

# Setup aliases for drush and composer. 
alias composer="sudo -u $SYSTEM_USER /usr/local/bin/composer"
alias drush="sudo -u $SYSTEM_USER $COMPOSER_DIR/vendor/bin/drush"

# Because of a bug in Drupal we constantly have to reset the permissions of the site directory. 
# See https://www.drupal.org/project/drupal/issues/3091285. 
function drupal_sites_permission_workaround() {
    chmod -R u+w "$WEB_DIR/sites/"
}

# Apache configuration paths
APACHE_CONFIG_SITE_AVAILABLE="/etc/apache2/sites-available/${INSTANCE_DOMAIN}.conf"
APACHE_CONFIG_SITE_ENABLED="/etc/apache2/sites-enabled/${INSTANCE_DOMAIN}.conf"