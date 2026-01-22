<?php

use \Drupal\Core\Site\Settings;
use \Drupal\Core\DependencyInjection\ContainerNotInitializedException;

/** gets a setting from 'settings.php' */
function get_setting($name) {
    return Settings::get($name);
}

/** sets a setting in 'settings.php' */
function set_setting(string $name, mixed $value): bool {
    $settings["settings"][$name] = (object)[
        "value" => $value,
        "required" => TRUE,
    ];

    return do_set_setting($settings);
}

/** sets the default database connection URL in settings.php */
function set_default_db_connection(string $url): bool {
    // Parse DB URL.
    $parts = parse_url($url);
    if ($parts === false) {
        throw new Exception("Invalid database URL.");
    }

    $scheme = strtolower($parts['scheme'] ?? '');
    if ($scheme !== 'mysql' && $scheme !== 'mariadb') {
        throw new Exception("Unsupported scheme '{$scheme}'. Use mysql:// or mariadb://");
    }

    $database = ltrim($parts['path'] ?? '', '/');
    if ($database === '') {
        throw new Exception("Database name missing in URL path (e.g. /dbname).");
    }

    // Build Drupal DB array (mysql driver).
    $db = [
        'driver' => 'mysql',
        'database' => $database,
        'username' => $parts['user'] ?? '',
        'password' => $parts['pass'] ?? '',
        'host' => $parts['host'] ?? '127.0.0.1',
        'port' => (string)($parts['port'] ?? '3306'),
        'prefix' => '',
        'namespace' => 'Drupal\\Core\\Database\\Driver\\mysql',
    ];

    // Prepare rewrite structure for SettingsEditor::rewrite().
    // IMPORTANT: This targets $databases['default']['default'] in settings.php.
    $settings = [];
    $settings['databases']['default']['default'] = (object) [
        'value' => $db,
        'required' => true,
    ];

    return do_set_setting($settings);
}

function do_set_setting(array $settings): bool {
    // find settings.php
    try {
        $filename = DRUPAL_ROOT . "/" . \Drupal::getContainer()->getParameter("site.path") . "/settings.php";
    } catch(ContainerNotInitializedException $t) {
        // Fallback to sites/default/settings.php
        $filename = DRUPAL_ROOT . '/sites/default/settings.php';
    }

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

    $ok = TRUE;
    if(class_exists('\Drupal\Core\Site\SettingsEditor')) {
        try {
            \Drupal\Core\Site\SettingsEditor::rewrite($filename, $settings);
        } catch(Throwable $t) {
            $ok = FALSE;
        }
    } elseif (function_exists('drupal_rewrite_settings')) {
        try {
            drupal_rewrite_settings($settings, $filename);
        } catch(Throwable $t) {
            $ok = FALSE;
        }
    } else {
        throw new Exception("Unsupported Drupal Release (need 'SettingsEditor' or 'drupal_rewrite_settings')");
    }

    // reset the file mode
    return chmod($filename, $old) && $ok;
}

/** Sets the trusted host to the specified domain */
function set_trusted_domain(string $domain): bool {
    return set_setting("trusted_host_patterns", [preg_quote($domain)]);
}

/** Sets up including a settings.php file from the given path */
function install_settings_include(array $paths): bool {
    // find the original filename
    $filename = DRUPAL_ROOT . "/" . \Drupal::getContainer()->getParameter("site.path") . "/settings.php";
    
    // read the original file
    $original_content = file_get_contents($filename);
    if ($original_content === FALSE) {
        return FALSE;
    }

    // remove any old <distillery-settings-includes>
    $pattern = '/\/\/(\s*)<distillery-settings-include>(.*?)\/\/(\s*)<\/distillery-settings-include>/s';
    $new_content = preg_replace($pattern, '', $originalContent);
    
    $code = "// <distillery-settings-include>>\n//\n// DO NOT MODIFY THIS BLOCK AND KEEP IT AT THE END OF THE FILE.\n// DO NOT REMOVE CONFIG TAGS\n";
    foreach ($paths as $path) {
        // escape the path to be included
        $the_path = "'" . addslashes($path) . "'";
        // resolve it (if it isn't absolute)
        if (!str_starts_with($path, '/')) {
            $the_path = '$app_root . \'/\' . $site_path . \'/\' . ' . $the_path;
        }

        // add code to include the file if it exists
        $code = $code . 'if (file_exists(' . $the_path . ')) { include_once ' . $the_path . '; }' . "\n";
    }
    $code = $code . "// </distillery-settings-include>\n";

    // and store the settings
    try {
        file_put_contents($filename, $original_content . $code);
    }  catch(Throwable $t) {
        return FALSE;
    }

    return TRUE;
}