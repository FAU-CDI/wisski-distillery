<?php

use \Drupal\Core\Site\Settings;

/** gets a setting from 'settings.php' */
function get_setting($name) {
    return Settings::get($name);
}

/** sets a setting in 'settings.php' */
function set_setting(string $name, mixed $value): bool {
    // find settings.php
    $filename = DRUPAL_ROOT . "/" . \Drupal::service("site.path") . "/settings.php";

    // setup user write permissions for the file
    $old = fileperms($filename);
    if ($old === FALSE) {
        return FALSE;
    }

    $new = 0777; // set all permissions
    if (!chmod($filename, $new)) {
        return FALSE;
    }
    
    // load install.inc
    if(is_file(DRUPAL_ROOT . "/internal/")) {
        include_once DRUPAL_ROOT . "/internal/core/includes/install.inc";
    } else {
        include_once DRUPAL_ROOT . "/core/includes/install.inc";
    }
    
    // update the provided setting
    $settings["settings"][$name] = (object)[
        "value" => $value,
        "required" => TRUE,
    ];

    // do the rewrite
    try {
        drupal_rewrite_settings($settings, $filename);
    } catch(Throwable $t) {
        throw $t; // DEBUG
        return FALSE;
    }


    // reset the file mode
    return chmod($filename, $old);
}

/** Sets the trusted host to the specified domain */
function set_trusted_domain(string $domain): bool {
    return set_setting("trusted_host_patterns", [preg_quote($domain)]);
}