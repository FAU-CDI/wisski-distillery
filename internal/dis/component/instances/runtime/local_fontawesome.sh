#!/bin/bash
set -e

FAINFO=`drush pm-list --pipe --type=module --filter=id=fontawesome`
if [[ -z "$FAINFO" ]]; then
    echo "Font Awesome is not installed, aborting"
    exit 0
fi

# make a temporary directory and cd into it
TEMPDIR=$(mktemp -d)
pushd "$TEMPDIR"
trap 'popd && rm -rf $TEMPDIR' EXIT

# curl the colorbox zip and unpack it
curl -L https://use.fontawesome.com/releases/v6.4.0/fontawesome-free-6.4.0-web.zip --output fontawesome.zip
unzip fontawesome.zip

# Prepare the fontawesome directory
chmod u+rw /var/www/data/project/web/
mkdir -p /var/www/data/project/web/libraries
rm -rf /var/www/data/project/web/libraries/fontawesome

# Move over the fontawesome zip file
mv fontawesome-* /var/www/data/project/web/libraries/fontawesome

# Update drush config to use local fontawesome
drush config:set --yes --input-format=yaml fontawesome.settings use_cdn false
drush config:set --yes --input-format=yaml fontawesome.settings use_shim false
drush cr