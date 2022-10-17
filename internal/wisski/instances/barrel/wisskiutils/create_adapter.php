<?php

/**
 * This script will automatically create a WissKI Salz Adapter for use within the distillery.
 * It will not update any existing adapter and is rather primitive.
 */

$argc = $_SERVER['argc']-3;
$argv = array_slice($_SERVER['argv'], 3);

// read parameters from the command line
if ($argc != 3) {
    die("Usage: drush php:script create_adapter.php INSTANCE_DOMAIN GRAPHDB_REPO HEADER");
}
$INSTANCE_DOMAIN = $argv[0];
$GRAPHDB_REPO = $argv[1];
$HEADER = $argv[2];

//
// PROPERTIES FOR THE ADAPTER
//

$id = 'default'; // id
$type = 'sparql11_with_pb'; // plugin
$machine_name = 'default'; // machine-name
$label = 'Default WissKI Distillery Adapter';
$description = 'Adapter for ' . $INSTANCE_DOMAIN; // description
$writable = TRUE; // writable
$is_preferred_local_store = TRUE; // is_preferred_local_store
$header = $HEADER; // header
$read_url = 'http://triplestore:7200/repositories/' . $GRAPHDB_REPO; // read_url
$write_url = 'http://triplestore:7200/repositories/' . $GRAPHDB_REPO . '/statements'; // write_url
$is_federatable = TRUE; // is_federatable
$default_graph_uri = 'https://' . $INSTANCE_DOMAIN . '/';
$same_as_properties = ['http://www.w3.org/2002/07/owl#sameAs']; // same_as_properties
$ontology_graphs = []; // ontology_graphs

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