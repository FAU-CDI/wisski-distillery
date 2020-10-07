#!/bin/bash

# This utility script can be used to configure the trusted host settings inside of settings.php. 
# It doesn't take care of corner cases and should only be used when needed

INSTANCE_DOMAIN="$(hostname -f)"
chmod u+w web/sites/default/settings.php
echo "" >> web/sites/default/settings.php
echo "\$settings['trusted_host_patterns'] = ['${INSTANCE_DOMAIN//\./\\\.}'];" >> web/sites/default/settings.php
chmod u-w web/sites/default/settings.php