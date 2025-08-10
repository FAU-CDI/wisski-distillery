<?php

use Drupal\Core\Url;
use Drupal\user\Entity\User;

/** lists all the users */
function list_users(): mixed {
    $users = [];
    foreach (User::loadMultiple(NULL) as $user) {
        $fields = array_map(function ($field) {
            return $field->getString();
        }, $user->getFields());
        if (empty($fields['name'])) continue;
        $users[] = $fields;
    }
    return $users;
}

function set_user_password($name, $password): bool {
    $user = user_load_by_name($name);
    if (!$user) return false;
    $user->setPassword($password);
    $user->save();
    return true;
}

function get_password_hash($name): string {
    $user = user_load_by_name($name);
    if (!$user) return "";
    return $user->get('pass')->getString();
}


function check_password_hash($password, $hash): bool {
    return \Drupal::service('password')->check($password, $hash);
}

function get_login_link(string $name, string $destination = "", bool $update_user = FALSE, bool $grant_admin = FALSE): string {
    $account = user_load_by_name($name);
    if (!$account) {
        if (!$update_user) return "";
        $account = create_new_disabled_user($name);
    }

    if ($update_user && $grant_admin) {
        $account->addRole('administrator');
    }

    if ($update_user) {
        $account->save();
    }

    return do_login_link($account, $destination);
}

function get_root_login_link(string $destination = ""): string {
    $uids = \Drupal::entityQuery('user')
        ->condition('status', 1)
        ->condition('roles', 'administrator')
        ->sort('uid', 'ASC')
        ->range(0, 1)
        ->accessCheck(FALSE)
        ->execute();

    if (empty($uids)) {
        throw new \Exception("No root user");
    };
    $root = User::load(reset($uids));
    if (empty($uids)) {
        throw new \Exception("failed to load root user");
    };
    
    return do_login_link($root, $destination);
}

function do_login_link(User $account, string $destination = ""): string {
    $timestamp = \Drupal::time()->getRequestTime();
    return Url::fromRoute(
        'user.reset.login',
        [
            'uid' => $account->id(),
            'timestamp' => $timestamp,
            'hash' => user_pass_rehash($account, $timestamp),
        ],
        [
            'absolute' => false,
            'query' => ['destination' => $destination],
            'language' => \Drupal::languageManager()->getLanguage($account->getPreferredLangcode()),
        ]
    )->toString();
}

function create_new_disabled_user($name) {
    $account = User::create([
        'name' => $name,
    ]);
    $account->activate();
    $account->save();
    return $account;
}