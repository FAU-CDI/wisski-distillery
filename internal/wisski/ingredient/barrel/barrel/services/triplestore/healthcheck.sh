#!/bin/bash
set -e

RDF4J_API_BASE="http://localhost:8080/rdf4j-server"

# check that the server itself is running
curl --silent --output /dev/null --fail "${RDF4J_API_BASE}/protocol";

# check for the default repository
if [[ -z "$RDF4J_REPOSITORY" ]]; then
    exit 0;
fi

# check that the default repository exists
curl --fail --silent --output /dev/null "${RDF4J_API_BASE}/repositories/${RDF4J_REPOSITORY}/size"