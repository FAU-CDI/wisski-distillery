#!/bin/bash
set -e

# Because we want to run graphdb as a limited user
# we need to make sure that the volumes are writable. 
# Because of that, we 'chown'

chown graphdb:graphdb /opt/graphdb/data
chown graphdb:graphdb /opt/graphdb/work
chown graphdb:graphdb /opt/graphdb/logs

# switch to the graphdb user, and run graphdb
su graphdb -c "/opt/graphdb/bin/graphdb $@"