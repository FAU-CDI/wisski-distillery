#!/bin/bash
set -e

# TODO: Delete system user

# read the lib/shared.sh and read the slug argument. 
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"
require_slug_argument

# Delete the apache configurationf files first. 
# This prevents drupal from being served. 
log_info " => Removing apache configuration files"
rm "$APACHE_CONFIG_SITE_ENABLED" || true
rm "$APACHE_CONFIG_SITE_AVAILABLE" || true

# Reload apache to apply the configuration. 
log_info " => Reloading apache"
systemctl reload apache2

# Delete the MySQL database next. 
log_info " => Deleting MySQL database '$MYSQL_DATABASE' and user '$MYSQL_USER'. "
mysql -e "DROP DATABASE IF EXISTS \`${MYSQL_DATABASE}\`;" || true
mysql -e "DROP USER IF EXISTS \`${MYSQL_USER}\`@localhost;"  || true

# Clear the GraphDB repository. 
log_info " => Deleting GraphDB repository '$GRAPHDB_REPO'"
curl -X DELETE http://127.0.0.1:7200/rest/repositories/$GRAPHDB_REPO/

log_info " => Deleting system user and group '$SYSTEM_USER'"
deluser "$SYSTEM_USER" || true
delgroup "$SYSTEM_USER" || true

# Finally remove any trace of the repository by removing the base directory. 
log_info " => Removing directory '$BASE_DIR'"
rm -rf "$BASE_DIR"

log_info " => Finished, '$INSTANCE_DOMAIN' has been removed. "