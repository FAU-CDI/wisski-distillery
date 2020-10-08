#!/bin/sh

# This script can be used to repatch EasyRDF when needed. 
cd /var/www/data/project/web/modules/contrib/wisski || exit 1
EASYRDF_RESPONSE="./vendor/easyrdf/easyrdf/lib/EasyRdf/Http/Response.php"
patch -N "$EASYRDF_RESPONSE" < "/patch/easyrdf.patch"