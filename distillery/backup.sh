#!/bin/bash
set -e

# read the lib/shared.sh and read the slug argument. 
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"


log_info " => Starting backup process. This might take a while. "
wait_for_sql

BACKUP_SLUG="$(date +%Y%m%dT%H%M%S)-$(randompw)"
BACKUP_INSTANCE_DIR="$DEPLOY_BACKUP_INPROGRESS_DIR/$BACKUP_SLUG"
BACKUP_FINAL_FILE="$DEPLOY_BACKUP_FINAL_DIR/$BACKUP_SLUG.tar.gz"

BACKUP_SQL_FILE="$BACKUP_INSTANCE_DIR/backup.sql"

BACKUP_TRIPLESTORE_DIR="$BACKUP_INSTANCE_DIR/triplestore"
BACKUP_TRIPLESTORE_SYSTEM="$BACKUP_TRIPLESTORE_DIR/system.nq"

BACKUP_FILESYSTEM_DIR="$BACKUP_INSTANCE_DIR/instances"

# create the backup directories
log_info " => Making '$BACKUP_INSTANCE_DIR'"
mkdir -p "$BACKUP_INSTANCE_DIR"
mkdir -p "$DEPLOY_BACKUP_FINAL_DIR"

# backup the configuration
cp "$CONFIG_FILE" "$BACKUP_INSTANCE_DIR/.env"

# Backup the sql backup
log_info " => Backing up the SQL database"
dockerized_mysqldump --all-databases > "$BACKUP_SQL_FILE"

# Backup the triplestore system
log_info " => Backing up Triplestore System"
mkdir -p "$BACKUP_TRIPLESTORE_DIR"
curl -X GET -H "Accept:application/n-quads" "http://localhost:7200/repositories/SYSTEM/statements?infer=false" > "$BACKUP_TRIPLESTORE_SYSTEM"

# backup individual repos
for REPO in `grep -oP '(?<=#repositoryID> ")[^"]+' $BACKUP_TRIPLESTORE_SYSTEM`; do
    log_info " => Backing up Triplestore Repository '$REPO'"
	curl -X GET -H "Accept:application/n-quads" "http://localhost:7200/repositories/$REPO/statements?infer=false" > "$BACKUP_TRIPLESTORE_DIR/repo_$REPO.nq"
done

# backup the filesystem
log_info " => Backing up instance filesystem"
cp -rpT "$DEPLOY_INSTANCES_DIR" "$BACKUP_FILESYSTEM_DIR"

# Package the backup into a .tar.gz
log_info " => Packaging '$BACKUP_FINAL_FILE'"
pushd "$BACKUP_INSTANCE_DIR" > /dev/null
tar --totals --checkpoint=10000 -zcf "$BACKUP_FINAL_FILE" .
popd > /dev/null

# Clean up the unpacked backup
log_info " => Cleaning up '$BACKUP_INSTANCE_DIR'"
rm -rf "$BACKUP_INSTANCE_DIR"

log_info " => Removing backups older than $MAX_BACKUP_AGE days"
find "$DEPLOY_BACKUP_FINAL_DIR" -type f -mtime "+$MAX_BACKUP_AGE" -print -exec rm -f {} \;