<?php

/**
 * This script will list all the URIs that this system is aware of.
 * This works by listing all the default graph uris of all the adapters.
 */

// iterate over all adapters
$storage = \Drupal::entityTypeManager()->getStorage('wisski_salz_adapter');
foreach ($storage->loadMultiple() as $adapter) {
    // read the configuration, and check if we have a default graph
    $conf = $adapter->getEngine()->getConfiguration();
    if(!array_key_exists('default_graph', $conf)) {
        continue;
    }

    // and echo it out
    echo $conf['default_graph'] . "\n";
}