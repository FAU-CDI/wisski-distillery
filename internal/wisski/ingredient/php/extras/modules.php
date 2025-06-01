<?php


function build_extended_infos(): array {
    $COMPOSER='composer';
    $DRUSH='drush';

    $root = drush_root($DRUSH);
    $drush_infos = drush_infos($DRUSH, $root);
    $composer_infos = composer_infos($COMPOSER, normpath($root, '..'));

    return extend_drush_infos($root, $drush_infos, $composer_infos);
}


/**
 * adds name and version of composer modules into the modules array
 */
function composer_infos(string $COMPOSER, string|null $pwd = null): array {
    $composer_name_to_path = array();
    foreach(exec_stdout_json("$COMPOSER info --path --format=json", $pwd)["installed"] as $module) {
        $name = $module['name'];
        $path = normpath($module['path']);
        $composer_name_to_path[$name] = $path;
    }

    $infos = array();
    foreach(exec_stdout_json("$COMPOSER info --format=json", $pwd)["installed"] as $module) {
        $name = $module['name'];
        $path = $composer_name_to_path[$name];
        $infos[$path] = array(
            'name' => $name,
            'path' => $path,
            'version' => $module['version'],
        );
    }
    return $infos;
}

function drush_root(string $DRUSH): string {
   return exec_stdout_json("$DRUSH status --format=json")['root']; 
}
function drush_infos(string $DRUSH, string $root): array {
    $infos = array();
    foreach(exec_stdout_json("$DRUSH pm:list --format=json") as $name => $module) {
        $path = normpath($root, $module['path']);
        $info = array(
            'name' => $module['name'],
            'display_name' => clean_display_name($module['display_name']),

            'path' => $path,
            'type' => $module['type'],

            'enabled' => $module['status'] === 'Enabled',
            'version' => $module['version'],
        );
        $infos[$path] = $info;
    }
    return $infos;
}

function extend_drush_infos(string $root, array $drush_infos, array $composer_infos): array {
    $infos = array();
    foreach($drush_infos as $path => $module) {
        $extended = unserialize(serialize($module));
        $extended['composer'] = get_composer_module($path, $root, $composer_infos);
        $infos[] = $extended;
    }
    return $infos;
}

/**
 * gets the composer module that matches $path
 */
function get_composer_module(string $path, string $root, array $composer_infos): array|null {
    while(true) {
        if (array_key_exists($path, $composer_infos)) {
            return $composer_infos[$path];
        }

        $parent = normpath($path, '..');
        if ($parent === null || $parent === $path || $parent === $root) {
            return null;
        }
        $path = $parent;
    }
    die("never reached");
}


function clean_display_name(string|null $name): string|null {
    if($name === null) {
        return null;
    }

    $index = strrpos($name, ' (');
    if ($index === false) {
        return $name;
    }
    return substr($name, 0, $index);
}


/**
 * Executes a command, and returns parsed standardard output json
 */
function exec_stdout_json(string $command, string|null $cwd = null): mixed {
    $stdout = exec_stdout_text($command, $cwd);

    $result = @json_decode($stdout, true);
    if ($result === null && json_last_error() !== JSON_ERROR_NONE) {
        throw new Exception("command '" . $command . "': failed to decode json: " . json_last_error_msg());
    }
    return $result;
}

/**
 * Executes a command, and returns standard output as a string.
 */
function exec_stdout_text(string $command, string|null $cwd = null): string {
    $result = exec_pipes($command, $cwd, null);
    if ($result === null) {
        throw new Exception("command '" . $command . "': exec failed");
    }

    $code = $result['code'];
    if ($code !== 0) {
        throw new Exception("command '" . $command . "': exec returned code $code");
    }

    $stdout = $result['stdout'];
    if (!is_string($stdout)) {
        throw new Exception("command '" . $command . "': failed to capture stdout");
    }
    
    return $stdout;
}

/**
 * Executes command in $cwd, passing in standard input and returning separate error and output text.
 */
function exec_pipes(string $command, string|null $cwd = null, string|null $stdin = null): array|null {
    $pipes = array();
    $process = proc_open($command, [
        0 => ['pipe', 'r'],
        1 => ['pipe', 'w'],
        2 => ['pipe', 'w'],
    ], $pipes, $cwd);

    if (!is_resource($process)) {
        return NULL;
    }

    if (is_string($stdin)) {
        fwrite($pipes[0], $stdin);
    }
    fclose($pipes[0]);

    $stdout = stream_get_contents($pipes[1]);
    fclose($pipes[1]);

    $stderr = stream_get_contents($pipes[2]);
    fclose($pipes[2]);

    return [
        'stdout' => $stdout,
        'stderr' => $stderr,
        'code' => proc_close($process),
    ];
}


/**
 * joins all the input paths, and normalizes the result
 */
function normpath(string|null ...$paths): string|null {
    $to_join = array();
    foreach($paths as $path) {
        if (!is_string($path)) {
            continue;
        }
        $to_join[] = $path;
    }
    if (count($to_join) === 0) {
        return null;
    }
    $path = realpath(join('/', $to_join));
    if (!is_string($path)) {
        return null;
    }
    return $path;
}