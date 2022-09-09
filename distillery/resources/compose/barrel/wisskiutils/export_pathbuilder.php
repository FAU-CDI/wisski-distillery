<?php

/**
 * This script will list all the URIs that this system is aware of.
 * This works by listing all the default graph uris of all the adapters.
 */

use Drupal\wisski_pathbuilder\Entity\WisskiPathEntity;

// load all the pathbuilders
$pbs = \Drupal::entityTypeManager()->getStorage('wisski_pathbuilder')->loadMultiple();

// map over the pathbuilders
$xmls = array_map(function($pb) {
    $xml = new \SimpleXMLElement("<pathbuilderinterface></pathbuilderinterface>");

    $paths = $pb->getAllPaths();
    foreach ($paths as $key => $path) {
        $id = $path->getID();

        $path = $pb->getPbPath($id);

        $pathChild = $xml->addChild("path");
        $pathObject = WisskiPathEntity::load($id);

        foreach ($path as $subkey => $value) {

            if (in_array($subkey, ['relativepath'])) {
                continue;
            }

            if ($subkey == "parent") {
                $subkey = "group_id";
            }

            $pathChild->addChild($subkey, htmlspecialchars($value));
        }

        $pathArray = $pathChild->addChild('path_array');
        foreach ($pathObject->getPathArray() as $subkey => $value) {
            $pathArray->addChild($subkey % 2 == 0 ? 'x' : 'y', $value);
        }

        $pathChild->addChild('datatype_property', htmlspecialchars($pathObject->getDatatypeProperty()));
        $pathChild->addChild('short_name', htmlspecialchars($pathObject->getShortName()));
        $pathChild->addChild('disamb', htmlspecialchars($pathObject->getDisamb()));
        $pathChild->addChild('description', htmlspecialchars($pathObject->getDescription()));
        $pathChild->addChild('uuid', htmlspecialchars($pathObject->uuid()));
        if ($pathObject->getType() == "Group" || $pathObject->getType() == "Smartgroup") {
            $pathChild->addChild('is_group', "1");
        } else {
            $pathChild->addChild('is_group', "0");
        }
        $pathChild->addChild('name', htmlspecialchars($pathObject->getName()));
    }

    // turn it into XML
    $dom = dom_import_simplexml($xml)->ownerDocument;
    $dom->formatOutput = TRUE;
    return $dom->saveXML();
}, $pbs);

echo json_encode($xmls);