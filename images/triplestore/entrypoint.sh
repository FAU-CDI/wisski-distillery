#!/bin/bash
set -e

# chown the volumes to graphdb
chown -R graphdb:graphdb /opt/graphdb/data
chown -R graphdb:graphdb /opt/graphdb/work
chown -R graphdb:graphdb /opt/graphdb/logs

# run graphdb as a limited user
su graphdb -c "/opt/graphdb/bin/graphdb $@"