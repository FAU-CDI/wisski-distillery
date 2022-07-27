#!/bin/bash
set -e

# read the lib/shared.sh
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"

# Read the 'GRAPHDB_ZIP' argument from the command line. 
# If it's not set, throw an error. 
GRAPHDB_ZIP=$1
if [ -z "$GRAPHDB_ZIP" ]; then
    log_error "Usage: system_install.sh GRAPHDB_ZIP"
    exit 1;
fi;


# print some general info on the screen
log_info "=> Preparing system to become a WissKI Distillery"
echo "This script will install or upgrade this system to become a WissKI distillery. "
echo "It is idempotent and can safely be run multiple times. "
sleep 5


# Install default system upgrades. 
log_info "=> Installing system updates"
apt-get update
apt-get upgrade -y

# install docker dependencies. 
log_info "=> Installing docker installer dependencies"
apt-get update
apt-get install -y curl

# install docker using an automated script. 
log_info "=> Installing docker"
curl -fsSL https://get.docker.com -o - | /bin/sh

# install docker-compose dependencies. 
log_info "=> Install docker-compose installer dependencies"
apt-get update
apt-get install -y python3-pip libffi-dev

# install docker-compose. 
log_info "=> Installing docker-compose"
pip3 install --upgrade docker-compose

log_info "=> Creating docker-compose directories and files"
mkdir -p "$DEPLOY_INSTANCES_DIR"
mkdir -p "$DEPLOY_WEB_DIR"
mkdir -p "$DEPLOY_SELF_DIR"
mkdir -p "$DEPLOY_RESOLVER_DIR"
mkdir -p "$DEPLOY_SSH_DIR"
mkdir -p "$DEPLOY_TRIPLESTORE_DIR"
mkdir -p "$DEPLOY_SQL_DIR"
mkdir -p "$DEPLOY_BACKUP_INPROGRESS_DIR"
mkdir -p "$DEPLOY_BACKUP_FINAL_DIR"

log_info "=> Creating 'distillery' network"
docker network create distillery || true

log_info "=> Creating 'docker-compose' files for the 'web'. "
install_resource_dir "compose/web" "$DEPLOY_WEB_DIR"

log_info " => Writing 'web' configuration file"
load_template "docker-env/web" \
    "DEFAULT_HOST" "${DEFAULT_DOMAIN}" \
    > "$DEPLOY_WEB_DIR/.env"

log_info "=> Creating 'docker-compose' files for the 'self'. "
install_resource_dir "compose/self" "$DEPLOY_SELF_DIR"

log_info "=> Creating 'docker-compose' files for the 'resolver'. "
install_resource_dir "compose/resolver" "$DEPLOY_RESOLVER_DIR"
touch "$DEPLOY_PREFIX_CONFIG"

log_info "=> Creating 'docker-compose' files for the 'ssh'. "
install_resource_dir "compose/ssh" "$DEPLOY_SSH_DIR"

# setup the lesencrypt host for the default domain
if [ -n "$LETSENCRYPT_HOST" ]; then
    LETSENCRYPT_HOST="$SELF_DOMAIN_SPEC"
fi;

log_info " => Writing 'self' configuration file"
load_template "docker-env/self" \
    "VIRTUAL_HOST" "${SELF_DOMAIN_SPEC}" \
    "LETSENCRYPT_HOST" "${LETSENCRYPT_HOST}" \
    "LETSENCRYPT_EMAIL" "${LETSENCRYPT_EMAIL}" \
    "TARGET" "${SELF_REDIRECT}" \
    "OVERRIDES_FILE" "${SELF_OVERRIDES_FILE}" \
    > "$DEPLOY_SELF_DIR/.env"

log_info " => Writing 'resolver' configuration file"
load_template "docker-env/resolver" \
    "VIRTUAL_HOST" "${SELF_DOMAIN_SPEC}" \
    "LETSENCRYPT_HOST" "${LETSENCRYPT_HOST}" \
    "LETSENCRYPT_EMAIL" "${LETSENCRYPT_EMAIL}" \
    "PREFIX_FILE" "${DEPLOY_PREFIX_CONFIG}" \
    "DEFAULT_DOMAIN" "${DEFAULT_DOMAIN}" \
    "LEGACY_DOMAIN" "${SELF_EXTRA_DOMAINS}" \
    > "$DEPLOY_RESOLVER_DIR/.env"

# copy over the directory
log_info "=> Creating 'docker-compose' files for the 'triplestore'. "
install_resource_dir "compose/triplestore" "$DEPLOY_TRIPLESTORE_DIR"

# copy the graphdb.zip
echo "Writing \"$DEPLOY_TRIPLESTORE_DIR/graphdb.zip\""
cp "$GRAPHDB_ZIP" "$DEPLOY_TRIPLESTORE_DIR/graphdb.zip"

# create data (volume) location
mkdir -p "$DEPLOY_TRIPLESTORE_DIR/data/data/"
mkdir -p "$DEPLOY_TRIPLESTORE_DIR/data/work/"
mkdir -p "$DEPLOY_TRIPLESTORE_DIR/data/logs/"

# copy over the sql resource directory, then ensure the data diretory for sql exists. 
log_info "=> Creating 'docker-compose' files for the 'sql'. "
install_resource_dir "compose/sql" "$DEPLOY_SQL_DIR"
mkdir -p "$DEPLOY_SQL_DIR/data/"

# Run all the updates via system_update.sh
log_info " => Running 'system_update.sh'"
bash "$SCRIPT_DIR/system_update.sh"

log_info "=> Waiting for sql to come up"
wait_for_sql

log_info "=> Creating '$DISTILLERY_BOOKKEEPING_DATABASE' database and '$DISTILLERY_BOOKKEEPING_TABLE' table"
load_template "bookkeeping/create.sql" "DATABASE" "$DISTILLERY_BOOKKEEPING_DATABASE" "TABLE" "$DISTILLERY_BOOKKEEPING_TABLE" | \
    dockerized_mysql

log_info "=> System installation finished, ready to distill. "