#!/bin/bash
set -e

# make a temporary directory and cd into it
TEMPDIR=$(mktemp -d)
pushd "$TEMPDIR"

# curl the colorbox zip and unpack it
curl -L https://github.com/jackmoore/colorbox/archive/master.zip --output master.zip
unzip master.zip

# make the directory for libraries, and remove the old colorbox installation
chmod u+rw /var/www/data/project/web/sites/default/
mkdir -p /var/www/data/project/web/sites/default/libraries/
rm -rf /var/www/data/project/web/sites/default/libraries/colorbox

# copy over the new installation
mv colorbox-master/ /var/www/data/project/web/sites/default/libraries/colorbox

# cleanup
popd
rm -rf "$TEMPDIR"
