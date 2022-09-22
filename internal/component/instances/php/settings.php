<?php

/** gets a setting from 'settings.php' */
function get_setting($name) {
    use \Drupal\Core\Site\Settings;
    return Settings::get($name);
}

/** sets a setting in 'settings.php' */
function set_setting($name, $value) {
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

    // find the filename
    $filename = DRUPAL_ROOT . "/" . \Drupal::service("site.path") . "/settings.php";
    drupal_rewrite_settings($settings, $filename);

    return True;
}