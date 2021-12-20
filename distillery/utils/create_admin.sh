#!/bin/sh
set -e

# read user
USER=$1
if [ -z "$USER" ]; then
    echo "Usage: create_admin.sh USERNAME"
    exit 1
fi

# read password
echo "Enter Password for $USER:"
read -s PASS
echo "Enter the same password again:"
read -s PASS2

if [ "$PASS" != "$PASS2" ]; then
    echo "Passwords not equal"
    exit 1
fi;

# create the user and add the admin role
cd /var/www/data/project/
drush user:create "$USER" --password="$PASS"
drush user-add-role administrator "$USER"