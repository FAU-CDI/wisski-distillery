#!/bin/bash
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

log_info " => Creating local directory structure at '$INSTANCE_BASE_DIR'"
mkdir -p "$INSTANCE_BASE_DIR"
install_resource_dir "compose/reserve" "$INSTANCE_BASE_DIR"

log_info " => Writing configuration file"
load_template "docker-env/reserve" \
    "VIRTUAL_HOST" "${INSTANCE_DOMAIN}" \
    "SLUG" "${SLUG}" \
    "LETSENCRYPT_HOST" "${LETSENCRYPT_HOST}" \
    "LETSENCRYPT_EMAIL" "${LETSENCRYPT_EMAIL}" \
    > "$INSTANCE_BASE_DIR/.env"


log_info " => Running and building image"
cd "$INSTANCE_BASE_DIR"
docker-compose build --pull
docker-compose pull

log_info " => Starting container"
docker-compose up -d

log_info " => $INSTANCE_DOMAIN has been reserved"