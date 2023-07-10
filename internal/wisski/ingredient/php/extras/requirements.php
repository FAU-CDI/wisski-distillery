<?php

/**
 * Returns a well-typed array of all requirements.
 * Relative URLs will be placed in the given domain.
 */
function get_requirements(string $public = ""): array {
    $results = [];

    $managerService = \Drupal::service('system.manager');
    $rendererService = \Drupal::service('renderer');

    $requirements = $managerService->listRequirements();
    foreach($requirements as $id => $req) {
        $title = $req['title'] ?? NULL;
        $weight = $req['weight'] ?? NULL;
        $severity = $req['severity'] ?? NULL;
        $value = $req['value'] ?? NULL;
        $description = $req['description'] ?? NULL;

        $results[] = array(
            "id" => $id,
            "title" => strip_tags(ensure_html($rendererService, $title)),
            "weight" => is_numeric($weight) ? $weight : 0,
            "severity" => is_numeric($severity) ? $severity : 0, 
            "value" => clean_html(ensure_html($rendererService, $value), $public),
            "description" => clean_html(ensure_html($rendererService, $description), $public),
        );
    }

    return $results;
}


/**
 * Ensures that the passed data is drupal html code.
 *
 * @param mixed $data
 * @return string
 */
function ensure_html(mixed $renderer, mixed $data): string {
  // already a string => return it!
  if (is_string($data)) {
	  return $data;
  }

  // it is null => return the empty string
  if ($data === NULL) {
	  return "";
  }

  // create a render array and render it!
  $rary = is_array($data) ? $data : ["#markup" => $data];
  return $renderer->renderPlain($rary);
}

/**
 * Parses source as an html fragment.
 * It then iterates over all '<a>' elements and performs the following modifications:
 * - Setup the distillery redirector (/next/?next=...) on relative links according to $public
 * - Make all links target="_blank" rel="noopener"
  */
function clean_html(string $source, string $public): string|bool {
    // fast path: we don't have a source
    if ($source === "" ) {
        return $source;
    }

    // trim the trailing '/' from the public URL.
    // this technically trims the beginning as well, but that should not be a problem
    $public = trim($public, '/');

    // attempt to parse a new document
    $doc = new DOMDocument();
    $ok = $doc->loadHTML('<?xml encoding="UTF-8">' . $source);
    if ($ok === FALSE) {
        return FALSE;
    }

    // replace all the <a href="..."> which are relative
    $as = $doc->getElementsByTagName('a');
    $modified = $as->count() > 0;
    foreach($as as $a) {
        // setup a link to open in a new tab
        $a->setAttribute('rel', 'noopener noreferer');
        $a->setAttribute('target', '_blank');

        // if we don't have a domain don't even bother replacing hrefs.
        if ($domain === "") {
            continue;
        }

        // take only href="/" relative to the current domain.
        $href = $a->getAttribute('href');
        if (!is_string($href) || !str_starts_with($href, '/')) {
            continue;
        }

        $a->setAttribute('href', $public . $href);
    }

    // we didn't modify the document => return as is
    // (no need to re-serialize)
    if (!$modified) {
        return $source;
    }

    // find the body
    $body = $doc->getElementsByTagName('body')->item(0);
    if ($body === NULL) {
        return FALSE;
    }

    // re-build the source document
    $source = "";

    // get the inner html§
    $nodes = $body->childNodes;
    foreach($nodes as $node) {
        $source .= $doc->saveHTML($node);
    }

    // and turn it back into html
    return $source;
}