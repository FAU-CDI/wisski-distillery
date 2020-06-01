#!/bin/bash
set -e

# read the lib/shared.sh and lib/slug.sh
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR"
source "$DIR/lib/lib.sh"
require_slug_argument

log_info " => Opening shell in '$COMPOSER_DIR'"

# cd into the right directory. 
cd "$COMPOSER_DIR"

# add /usr/local/bin (for composer) and the vendor bin (for drush) to path
# and open a bash shell as www-data there. 
sudo -u "$SYSTEM_USER" PATH="$COMPOSER_DIR/vendor/bin:/usr/local/bin:$PATH" /bin/bash