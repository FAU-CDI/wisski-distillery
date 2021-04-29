#!/bin/bash

# This utility script can be used to configure the trusted host settings inside of settings.php. 
# It doesn't take care of corner cases and should only be used when needed. 

INSTANCE_DOMAIN="$(hostname -f)"
INSTANCE_DOMAIN="${INSTANCE_DOMAIN%.wisski}"

TRUSTED_HOST_PATTERN="${INSTANCE_DOMAIN//\./\\\\.}"
TRUSTED_HOST_PATTERNS='["'$TRUSTED_HOST_PATTERN'"]'

echo "Setting 'trusted_host_patterns' to $TRUSTED_HOST_PATTERNS"
bash /utils/settings_php_set.sh 'trusted_host_patterns' "$TRUSTED_HOST_PATTERNS"