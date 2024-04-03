#!/bin/bash
set -e

# if the user is not www-data, re-invoke self as www-data
if [ "$USER" != "www-data" ]; then
    sudo -u www-data /bin/bash /user_shell.sh "$@"
    exit $?
fi

# Now start a shell in the proper path
cd "/var/www/data/project"
export "PATH=/var/www/data/project/vendor/bin:$PATH"

# Re-invoke the actual command
/bin/bash "$@"