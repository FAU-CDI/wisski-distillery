#!/bin/bash
set -e

# read the lib/shared.sh
DISABLE_LOG=1
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"

# wait for sql to come up
wait_for_sql > /dev/null
dockerized_mysql_interactive "$@"
