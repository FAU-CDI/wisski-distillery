#!/bin/bash
set -e

# read the lib/shared.sh
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"

log_info "=> Rebuilding and restarting 'web' stack"
update_stack "$DEPLOY_WEB_DIR"

log_info "=> Rebuilding and restarting 'self' stack"
update_stack "$DEPLOY_SELF_DIR"

# build and start the ssh server
log_info "=> Rebuilding and restarting 'ssh' stack"
update_stack "$DEPLOY_SSH_DIR"

# build and start the triplestore
log_info "=> Rebuilding and restarting 'triplestore' stack"
update_stack "$DEPLOY_TRIPLESTORE_DIR"

# build and start the triplestore
log_info "=> Rebuilding and restarting 'sql' stack"
update_stack "$DEPLOY_SQL_DIR"

log_info " => Updating Prefix Config"
cd "$DIR"
bash update_prefix_config.sh

log_info "=> Rebuilding and restarting 'resolver' stack"
update_stack "$DEPLOY_RESOLVER_DIR"

log_info "=> System up-to-date. "