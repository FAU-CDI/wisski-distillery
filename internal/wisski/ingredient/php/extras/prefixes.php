<?php

/**
 * list_prefixes lists all content prefixes known to this WissKI.
 * Prefixes are not filtered, and may contain duplicates.
 */ 
function list_prefixes() {
    $prefixes = [];
    $storage = \Drupal::entityTypeManager()->getStorage('wisski_salz_adapter');
    foreach ($storage->loadMultiple() as $adapter) {
        // load all the prefixes from the triplestore
        $engine = $adapter->getEngine();
        getTriplestorePrefixes($adapter->getEngine(), $prefixes);

        // read the configuration to check if we have a default graph
        $conf = $engine->getConfiguration();
        if(!array_key_exists('default_graph', $conf)) {
            continue;
        }
        $prefixes[] = $conf['default_graph'];
    }
    return $prefixes;
}

function getTriplestorePrefixes($engine, &$prefixes) {
    // some adapters don't support a query method!
    if (!method_exists($engine, 'directQuery')) return NULL;

    $results = $engine->directQuery('
    select distinct ?base where {
        {
            select distinct ?iri where {
                {
                    select distinct (?s as ?iri) { ?s ?p ?o  }
                } union {
                    select distinct (?o as ?iri) { ?s ?p ?o FILTER(isiri(?o)) }
                }
            }  
        }
        BIND(replace(str(?iri), "/[^/]*/?$", "/") as ?base)
        FILTER(!REGEX(?base, "/wisski/navigate/[\\\\d]+/$"))
    } ORDER BY ?base');
    if (!$results) return FALSE;

    foreach($results as $result) {
        $prefixes[] = $result->base->getValue();
    }

    return TRUE;
} 

