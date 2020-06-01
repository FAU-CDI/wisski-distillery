#!/bin/bash
set -e

# read the lib/shared.sh and lib/slug.sh
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"
require_slug_argument

# TODO: Figure out if this is enough. 
echo " => Running 'composer update'"
cd "$COMPOSER_DIR"
drupal_sites_permission_workaround
composer update