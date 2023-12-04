#!/bin/bash
set -e

echo "=> Setting up filesystem permissions"
chmod 777 /var/www/data/project/web/sites/default/
trap "chmod 755 /var/www/data/project/web/sites/default/" EXIT

echo "=> Creating 'sites/default/libraries/colorbox/' directory"
mkdir -p /var/www/data/project/web/sites/default/libraries/colorbox

echo "=> Downloading 'jquery.colorbox-min.js' and 'LICENSE.md'"
curl -L https://raw.githubusercontent.com/jackmoore/colorbox/master/LICENSE.md -o /var/www/data/project/web/sites/default/libraries/colorbox/LICENSE.md
curl -L https://raw.githubusercontent.com/jackmoore/colorbox/master/jquery.colorbox-min.js -o /var/www/data/project/web/sites/default/libraries/colorbox/jquery.colorbox-min.js

echo "=> Done"