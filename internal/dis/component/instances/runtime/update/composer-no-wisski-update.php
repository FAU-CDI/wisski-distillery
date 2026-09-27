<?php declare(strict_types=1);

/**
 * composer-update-no-wisski.php
 * 
 * Script to update the composer.json file used by the distillery, excluding the Wisski package.
 * 
 * Usage:
 *  php composer-update-no-wisski.php
 */

// ==================== CONFIG ====================
$config = [
    'excluded' => ['drupal/wisski'],
];

// ==================== READ PACKAGES FROM composer.json ====================
$file = getcwd() . '/composer.json';
$json = json_decode(file_get_contents($file), true);
$packages = array_merge(
    array_keys($json['require'] ?? []),
    array_keys($json['require-dev'] ?? [])
);
$update = array_values(array_diff($packages, $config['excluded']));

// ==================== BUILD COMMAND ====================
$cmd = 'composer update ' . implode(' ', array_map('escapeshellarg', $update));

foreach ($argv as $i => $arg) {
    if ($i === 0) continue;
    $cmd .= ' ' . escapeshellarg($arg);
}

// ==================== EXECUTE WITH STREAM FORWARDING ====================
$descriptorspec = [
    0 => STDIN,
    1 => STDOUT,
    2 => STDERR,
];

$process = proc_open($cmd, $descriptorspec, $pipes);
$returnCode = proc_close($process);

exit($returnCode);