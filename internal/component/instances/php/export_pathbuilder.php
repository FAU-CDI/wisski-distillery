<?php

use Drupal\wisski_pathbuilder\Entity\WisskiPathEntity;

/** all_xml lists all pathbuilders, and returns the corresponding xml */
function all_xml(): object {
    $all = \Drupal::entityTypeManager()->getStorage('wisski_pathbuilder')->loadMultiple();
    return (object)array_map("entity_to_xml", $all);
}


/** all_list lists the ids of all pathbuilders */
function all_list(): Array {
    return array_keys(\Drupal::entityQuery('wisski_pathbuilder')->execute());
}

/** one_xml serializes a single pathbuilder as xml */
function one_xml(string $id): string {
    $pb = \Drupal::entityTypeManager()->getStorage('wisski_pathbuilder')->load($id);
    if ($pb === NULL) {
        return "";
    }
    return entity_to_xml($pb);
}

// =================================================================================
// =================================================================================


function entity_to_xml($pb) {
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
}