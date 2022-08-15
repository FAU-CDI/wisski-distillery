#!/bin/bash
set -e

# read the lib/shared.sh
DISABLE_LOG=0
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"

# wait for sql to come up
wait_for_sql > /dev/null

echo "Creating new MySQL user with root privileges. "
read -p 'Enter Username:' MYSQL_USER
read -sp 'Enter password:' MYSQL_PASSWORD

if ! is_valid_slug "$MYSQL_USER"; then
    echo "Not a valid username: ${MYSQL_USER}"
    echo "User must be alphanumeric for sql injection reasons. "
    echo "You can always create a user manually. "
    exit 1
fi

if ! is_valid_slug "$MYSQL_PASSWORD"; then
    echo "Not a valid password: ${MYSQL_PASSWORD}"
    echo "Password must be alphanumeric for sql injection reasons. "
    echo "You can always create a user manually. "
    exit 1
fi

dockerized_mysql -e "CREATE USER \`${MYSQL_USER}\`@'%' IDENTIFIED BY '${MYSQL_PASSWORD}'; GRANT ALL PRIVILEGES ON *.* TO \`${MYSQL_USER}\`@\`%\` WITH GRANT OPTION; FLUSH PRIVILEGES;"

log_info "Created user ${MYSQL_USER}"