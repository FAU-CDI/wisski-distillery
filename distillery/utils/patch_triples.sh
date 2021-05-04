#!/bin/sh

# This script can be used to repatch EasyRDF when needed. 
cd /var/www/data/project/web/modules/contrib/wisski/ || exit 1
TRIPLESTABCONTROLLER="./wisski_adapter_sparql11_pb/src/Controller/Sparql11TriplesTabController.php"
patch -N "$TRIPLESTABCONTROLLER" < "/patch/triples.patch"