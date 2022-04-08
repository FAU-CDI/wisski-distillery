#!/bin/bash
set -e

# read the lib/shared.sh and read the slug argument. 
DISABLE_LOG=1
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"
unset DISABLE_LOG
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

NEW_DOMAIN="wisski.data.fau.de"
NEW_INSTANCE_DOMAIN="$SLUG.$NEW_DOMAIN"
NEW_INSTANCE_DOMAIN="$(echo "$NEW_INSTANCE_DOMAIN" | tr '[:upper:]' '[:lower:]')"
NEW_INSTANCE_BASE_DIR="$DEPLOY_INSTANCES_DIR/$NEW_INSTANCE_DOMAIN"
NEW_INSTANCE_DATA_DIR="$NEW_INSTANCE_BASE_DIR/data/"
NEW_LETSENCRYPT_HOST="$NEW_INSTANCE_DOMAIN"

CONFIG_FILE="$INSTANCE_BASE_DIR/.env"
NEW_CONFIG_FILE="$INSTANCE_BASE_DIR/.envnew"
OLD_CONFIG_FILE="$INSTANCE_BASE_DIR/.envold"

log_info " => New Configuration for \'$NEW_DOMAIN\'"

echo "NEW_DOMAIN=$NEW_DOMAIN"
echo "NEW_INSTANCE_DOMAIN=$NEW_INSTANCE_DOMAIN"
echo "NEW_INSTANCE_BASE_DIR=$NEW_INSTANCE_BASE_DIR"

echo " => Preparing new configuration file"

load_template "docker-env/barrel" \
    "REAL_PATH" "${NEW_INSTANCE_DATA_DIR}" \
    "GLOBAL_AUTHORIZED_KEYS_FILE" "${GLOBAL_AUTHORIZED_KEYS_FILE}" \
    "VIRTUAL_HOST" "${NEW_INSTANCE_DOMAIN}" \
    "SLUG" "${SLUG}" \
    "LETSENCRYPT_HOST" "${NEW_LETSENCRYPT_HOST}" \
    "LETSENCRYPT_EMAIL" "${LETSENCRYPT_EMAIL}" \
    "DISTILLERY_DIR" "${DIR}" | tee "$NEW_CONFIG_FILE"

while true; do
    log_info " => I'm about to make breaking changes. "
    read -p "This can not be undone. Please type 'y' to continue: " yn
    case $yn in
        [Yy]* ) break;;
        * ) echo "Abort. "; exit 1;;
    esac
done

log_info " => Shutting down old system"

cd "$INSTANCE_BASE_DIR"
docker-compose down
cd "$DIR"

log_info " => Writing new configuration files"
mv "$CONFIG_FILE" "$OLD_CONFIG_FILE"
mv "$NEW_CONFIG_FILE" "$CONFIG_FILE"

log_info " => Moving base directory"
mv "$INSTANCE_BASE_DIR" "$NEW_INSTANCE_BASE_DIR"

log_info " => Updating bookeeping"
dockerized_mysql -D "$DISTILLERY_BOOKKEEPING_DATABASE" -e "UPDATE \`$DISTILLERY_BOOKKEEPING_TABLE\` SET \`filesystem_base\`='$NEW_INSTANCE_BASE_DIR' WHERE \`slug\`='$SLUG';"

log_info " => Starting in new location"
cd "$NEW_INSTANCE_BASE_DIR"
docker-compose up -d

log_info " => Updating trusted hosts"
docker-compose exec barrel /user_shell.sh /utils/set_trusted_host.sh

log_info " => We should be moved to '$NEW_INSTANCE_DOMAIN' now"