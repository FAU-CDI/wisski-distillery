#!/bin/bash
set -e

# Ensure the RDF4J data directory is writable by the tomcat user
chown -R tomcat: /var/rdf4j

# Run entrypoint2.sh (start server, init repo, wait) as the tomcat user
exec runuser -u tomcat -- /bin/bash /entrypoint2.sh "$@"