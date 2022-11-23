<?php

use Drupal\user\Entity\User;

/** lists all the users */
function list_users() {

    $usernames = [];
    $users = User::loadMultiple(NULL);
    foreach($users as $user){
        $name = $user->get('name')->getString();
        if(empty($name)) continue;
        $usernames[] = $name;
    }
    return $usernames;
}
