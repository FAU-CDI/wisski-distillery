<?php declare(strict_types=1);

/**
 * update-composer-json.php
 * 
 * Script to update the composer.json file used by the distillery.
 * 
 * It will generally fix version constraints to a ~major.minor constraint; and update other dependencies via specific constraint handling.
 * 
 * By default, it simply writs the updated composer.json to stdout.
 * To write the updated composer.json to a file, use the --write option, this will automatically create a backup of the original composer.json file.
 * 
 * Usage:
 *  php update-composer-json.php [--write]
 */
$config = [
    // pins that aren't touched
    'untouched_pins' => [
        'drupal/wisski',
    ],

    // constraints to add to the given sections
    'add_require_constraints' => [
    ],
    'add_require-dev_constraints' => [
    ],


    // constraints to remove from the given sections
    'remove_constraints' => [
    ],

    // constraints that are written if they exist
    'custom_constraints' => [
        // Use composer/installers ~2.3.0
        "composer/installers" => "~2.3.0",

        // Use drupal core 11.4.x
        "drupal/core-composer-scaffold" => "~11.4.8",
        "drupal/core-project-message" => "~11.4.8",
        "drupal/core-recommended" => "~11.4.8",
    ],

    // composer file names
    'composer_json' => 'composer.json',
    'composer_lock' => 'composer.lock',
];

function loadJson(string $path): array {
    $data = json_decode(file_get_contents($path), true, 512, JSON_THROW_ON_ERROR);
    if (!is_array($data)) {
        throw new RuntimeException('Invalid JSON file: ' . $path);
    }
    return $data;
}

function isBeta(string $v): bool {
    return preg_match('/(?:alpha|beta|rc|dev)/i', $v) === 1;
}

function pinMinor(string $version): string {
    $version = ltrim($version, 'v');
    return '~' . $version;
}

$options = getopt('', ['write']);
$writeToFile = isset($options['w']) || isset($options['write']);

$lock = loadJson($config['composer_lock']);
$json = loadJson($config['composer_json']);

$locked = [];
foreach (($lock['packages'] ?? []) + ($lock['packages-dev'] ?? []) as $pkg) {
    $locked[$pkg['name']] = $pkg['version'];
}

foreach (['require', 'require-dev'] as $section) {
    // Iterate over constraints in the section
    foreach (($json[$section] ?? []) as $pkg => $constraint) {
        if (!isset($locked[$pkg])) continue;

        if (in_array($pkg, $config['untouched_pins'])) continue;


        if (isset($config['custom_constraints'][$pkg])) {
            $json[$section][$pkg] = $config['custom_constraints'][$pkg];
            continue;
        }

        if (isset($config['add_constraints'][$pkg])) {
            
            continue;
        }

        // Do the actual lock
        $ver = $locked[$pkg];
        if (isBeta($ver)) continue;

        $json[$section][$pkg] = pinMinor($ver);
    }

    // Add new constraints to the section
    foreach($config['add_' . $section . '_constraints'] as $pkg => $constraint) {
        $json[$section] ??= [];
        $json[$section][$pkg] = $constraint;
    }

    // remove the given constraints.
    foreach($config['remove_constraints'] as $pkg) {
        unset($json[$section][$pkg]);
    }
}

$output = json_encode($json, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_THROW_ON_ERROR) . "\n";

if ($writeToFile) {
    $timestamp = time();
    $backupFile = 'composer-backup-' . $timestamp . '.json';
    copy($config['composer_json'], $backupFile);
    file_put_contents($config['composer_json'], $output);
} else {
    echo $output;
}
?>