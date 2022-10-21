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


function entity_to_xml($pathbuilderEntity) {
    // NOTE: This function is verbatum copied from wisski_pathbuilder/src/PathbuilderManager.php.
    // The original code is licensed GPL-2-or-later, we choose GPL 3.0. 
    //
    // As per section 13 of GPL 3.0, we can reuse it under AGPL-3.0 (which this project is licensed under).
    
    // Create initial XML tree.
    $xmlTree = new \SimpleXMLElement("<pathbuilderinterface></pathbuilderinterface>");

    // Get the paths.
    $paths = $pathbuilderEntity->getPbPaths();

    // Iterate over every path.
    foreach ($paths as $key => $path) {
      $pathbuilder = $pathbuilderEntity->getPbPath($path['id']);
      $pathChild = $xmlTree->addChild("path");
      $pathObject = WisskiPathEntity::load($path['id']);

      foreach ($pathbuilder as $subkey => $value) {
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
      }
      else {
        $pathChild->addChild('is_group', "0");
      }
      $pathChild->addChild('name', htmlspecialchars($pathObject->getName()));

    }
    
    // Create XML DOM.
    $dom = dom_import_simplexml($xmlTree)->ownerDocument;
    $dom->formatOutput = TRUE;

    return $dom->saveXML();
}