#!/bin/bash
RDF4J_API_BASE="http://localhost:8080/rdf4j-server"

# spellchecker:words rdf4j spoc posc openrdf

# waits for rdf4j to be ready
function rdf4j_wait() {
    while true; do
        if curl --silent --output /dev/null --fail "${RDF4J_API_BASE}/protocol"; then
            break
        fi
        
        echo "RDF4J has not yet started ..."
        sleep 1
    done
}

# Checks if a repository already exists
function rdf4j_has_repository() {
    local RDF4J_REPOSITORY=$1
    curl --fail --silent --output /dev/null "${RDF4J_API_BASE}/repositories/${RDF4J_REPOSITORY}/size" || return 1
    return 0
}

# Creates a repository with the given name and label.
function rdf4j_create_repository() {
    local RDF4J_REPOSITORY=$1
    local RDF4J_LABEL=$2

    curl --fail --request PUT "${RDF4J_API_BASE}/repositories/${RDF4J_REPOSITORY}" -H 'Content-Type: text/turtle' --data-binary @- <<EOF
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#>.
@prefix config: <tag:rdf4j.org,2023:config/>.
[] a config:Repository ;
config:rep.id "${RDF4J_REPOSITORY}" ;
rdfs:label "${RDF4J_LABEL}" ;
config:rep.impl [
    config:rep.type "openrdf:SailRepository" ;
    config:sail.impl [
        config:sail.type "openrdf:NativeStore" ;
        config:native.tripleIndexes "spoc,posc"
    ]
].
EOF

    if [[ $? -ne 0 ]]; then
        return 1;
    fi

    return 0
}

# Main logic for initializing the container
function init_container() {
    if [[ -z "$RDF4J_REPOSITORY" ]]; then
        echo "[INIT SCRIPT]: 'RDF4J_REPOSITORY' is unset or empty.";
        echo "[INIT SCRIPT]: Skipping Repository Creation.";
        return 0;
    fi

    rdf4j_wait

    if rdf4j_has_repository "$RDF4J_REPOSITORY"; then
        echo "[INIT SCRIPT]: Repository '$RDF4J_REPOSITORY' already exists.";
        return 0;
    fi

    echo "[INIT SCRIPT]: Initializing repository '$RDF4J_REPOSITORY'";
    rdf4j_create_repository "$RDF4J_REPOSITORY" "$RDF4J_REPOSITORY_LABEL"

    if ! rdf4j_has_repository "$RDF4J_REPOSITORY"; then
        echo "[INIT SCRIPT]: Repository '$RDF4J_REPOSITORY' creation failed.";
        return 1;
    fi
    echo "[INIT SCRIPT]: Repository '$RDF4J_REPOSITORY' created.";
    return 0;
}

# start the original command 
"$@" &
CMD_PID=$!

# forward signals to the original entrypoint
trap 'kill -s SIGTERM $CMD_PID || exit 1' SIGTERM
trap 'kill -s SIGINT $CMD_PID || exit 1' SIGINT
trap 'kill -s SIGHUP $CMD_PID || exit 1' SIGHUP

# run the initialization logic
init_container || exit 1;

# wait for the original entrypoint to exit
wait
