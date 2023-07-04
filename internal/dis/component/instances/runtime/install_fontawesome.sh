#!/bin/bash
set -e

# make a temporary directory and cd into it
TEMPDIR=$(mktemp -d)
pushd "$TEMPDIR"

# curl the colorbox zip and unpack it
curl -L https://use.fontawesome.com/releases/v6.4.0/fontawesome-free-6.4.0-web.zip --output fontawesome.zip
unzip fontawesome.zip

# Prepare the fontawesome directory
chmod u+rw /var/www/data/project/web/
mkdir -p /var/www/data/project/web/libraries
rm -rf /var/www/data/project/web/libraries/fontawesome

# Move over the fontawesome zip file
mv fontawesome-* /var/www/data/project/web/libraries/fontawesome

# cleanup
popd
rm -rf "$TEMPDIR"
