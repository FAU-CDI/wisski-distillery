#!/bin/bash
set -e

# make the target directory and install
mkdir -p /var/www/data/project/web/sites/default/libraries/colorbox
curl -L https://raw.githubusercontent.com/jackmoore/colorbox/master/LICENSE.md -o /var/www/data/project/web/sites/default/libraries/colorbox/LICENSE.md
curl -L https://raw.githubusercontent.com/jackmoore/colorbox/master/jquery.colorbox-min.js -o /var/www/data/project/web/sites/default/libraries/colorbox/jquery.colorbox-min.js
