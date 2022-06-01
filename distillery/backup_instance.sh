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
    echo "I can't open a shell there. "
    exit 1
fi;

# Read everything from the database
read -r INSTANCE_BASE_DIR MYSQL_DATABASE MYSQL_USER GRAPHDB_REPO GRAPHDB_USER <<< "$(sql_bookkeep_load "${SLUG}" "filesystem_base,sql_database,sql_user,graphdb_repository,graphdb_user" | tail -n +2)"

# prepare the backup
log_info " => Preparing Backup Configuration"

BACKUP_SLUG="$SLUG-$(date +%Y%m%dT%H%M%S)-$(randompw)"
BACKUP_INSTANCE_DIR="$DEPLOY_BACKUP_INPROGRESS_DIR/$BACKUP_SLUG"

BACKUP_LOG_FILE="$BACKUP_INSTANCE_DIR/backup.log"
BACKUP_INFO_FILE="$BACKUP_INSTANCE_DIR/info.txt"
BACKUP_SQL_FILE="$BACKUP_INSTANCE_DIR/$MYSQL_DATABASE.sql"
BACKUP_NQ_FILE="$BACKUP_INSTANCE_DIR/$GRAPHDB_REPO.nq"
BACKUP_FS_DIR="$BACKUP_INSTANCE_DIR/$SLUG"
mkdir -p "$BACKUP_FS_DIR"

BACKUP_FINAL_FILE="$2"
if [ -z "$BACKUP_FINAL_FILE" ]; then
    BACKUP_FINAL_FILE="$DEPLOY_BACKUP_FINAL_DIR/$BACKUP_SLUG.tar.gz"
fi

echo "Destination: $BACKUP_FINAL_FILE"
if [ -z "$KEEPALIVE" ]; then
    echo "Keepalive: false (set KEEPALIVE variable for a consistent state)"
else
    echo "Keepalive: true (unset the KEEPALIVE variable for a consistent state)"
fi

BACKUP_START="$(date +%s)"

function do_the_backup() {
    # cd into the right directory
    cd "$INSTANCE_BASE_DIR"

    # stop
    if [ -z "$KEEPALIVE" ]; then
      log_info " => Shutting down running system"
      docker-compose down
    fi

    # system info
    log_info " => Backup up system information"
    /bin/bash "$DIR/info.sh" "$SLUG" > "$BACKUP_INFO_FILE"

    # database
    log_info " => Backing up MySQL database '$MYSQL_DATABASE'"
    dockerized_mysqldump --databases "$MYSQL_DATABASE" > "$BACKUP_SQL_FILE" || echo "Failed, continuing anyways ..."

    # triplestore
    log_info " => Backing up GraphDB repository '$GRAPHDB_REPO'"
    curl -X GET -H "Accept:application/n-quads" $GRAPHDB_AUTH_FLAGS "http://127.0.0.1:7200/repositories/${GRAPHDB_REPO}/statements?infer=false" > "$BACKUP_NQ_FILE" || echo "Failed, contiuing anyways ..."

    # filesystem
    log_info " => Backing up filesystem from '$INSTANCE_BASE_DIR'"
    cp -rpT "$INSTANCE_BASE_DIR" "$BACKUP_FS_DIR" || echo "Failed, continuing anyways ..."

    # restart
    if [ -z "$KEEPALIVE" ]; then
      log_info " => Starting up system"
      docker-compose up -d
    fi
}

# do the actual backup, writing it to a file
do_the_backup 2>&1 | tee "$BACKUP_LOG_FILE"

# list before packaging
log_info " => All backup files have been collected"
ls -alh "$BACKUP_INSTANCE_DIR"
du -hs "$BACKUP_INSTANCE_DIR"

# package up the backup
log_info " => Packaging \"$BACKUP_FINAL_FILE\" "
pushd "$BACKUP_INSTANCE_DIR" > /dev/null
tar --totals --checkpoint=10000 -zcf "$BACKUP_FINAL_FILE" .
popd > /dev/null

# delete the temporary directory
log_info " => Deleting temporary directories"
rm -r "$BACKUP_INSTANCE_DIR"

# and finish!
log_info " => Finished making backup"

DURATION=$[ $(date +%s) - ${BACKUP_START} ]
SIZE=$(wc -c < "$BACKUP_FINAL_FILE")

echo "Output file: '$BACKUP_FINAL_FILE'"
echo "Size:        $SIZE bytes"
echo "Duration:    $DURATION seconds"