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

function get_login_link($name): string {
    $account = user_load_by_name($name);
    if (!$account) return "";
    
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
            'query' => ['destination' => '/'],
            'language' => \Drupal::languageManager()->getLanguage($account->getPreferredLangcode()),
        ]
    )->toString();
}
