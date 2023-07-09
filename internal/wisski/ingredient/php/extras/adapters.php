<?php


/**
 * Creates an adapter for the distillery
 */
function create_distillery_adapter(string $LABEL, string $MACHINE_NAME, string $DESCRIPTION, string $INSTANCE_DOMAIN, string $GRAPHDB_REPO, string $GRAPHDB_USER, string $GRAPHDB_PASSWORD) {

    //
    // PROPERTIES FOR THE ADAPTER
    //
    $id = 'default'; // id
    $type = 'sparql11_with_pb'; // plugin
    $machine_name = $MACHINE_NAME; // machine-name
    $label = $LABEL;
    $description = $DESCRIPTION; // description
    $writable = TRUE; // writable
    $is_preferred_local_store = TRUE; // is_preferred_local_store
    $read_url = 'http://triplestore:7200/repositories/' . $GRAPHDB_REPO; // read_url
    $write_url = 'http://triplestore:7200/repositories/' . $GRAPHDB_REPO . '/statements'; // write_url
    $is_federatable = TRUE; // is_federatable
    $default_graph_uri = 'https://' . $INSTANCE_DOMAIN . '/';
    $same_as_properties = ['http://www.w3.org/2002/07/owl#sameAs']; // same_as_properties
    $ontology_graphs = []; // ontology_graphs

    // header
    $header = "";
    if ($GRAPHDB_USER !== "" && $GRAPHDB_PASSWORD !== "") {
        $header = $GRAPHDB_USER . ":" . $GRAPHDB_PASSWORD;
        $header = base64_encode($header);
    }
    
    //
    // Do the creation!
    //
    
    $storage = \Drupal::entityTypeManager()->getStorage('wisski_salz_adapter');
    $adapter = $storage->create([
        "id" => $id,
        "label" => $label,
        "description" => $description,
    ]);
    $adapter->setEngineConfig([
        "id" => $type,
        "machine-name" => $machine_name,
        "header" => $header,
        "writeable" => $writable,
        "is_preferred_local_store" => $is_preferred_local_store,
        "read_url" => $read_url,
        "write_url" => $write_url,
        "is_federatable" => $is_federatable,
        "default_graph" => $default_graph_uri,
        "same_as_properties" => $same_as_properties,
        "ontology_graphs" => $ontology_graphs,
    ]);
    $adapter->save();
}
