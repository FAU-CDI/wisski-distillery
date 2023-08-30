<?php

/**
 * Gets the actve theme of this WissKI
 *
 * @return string
 */
function get_active_theme(): string {
    $theme = \Drupal::service('theme.manager')->getActiveTheme();
    if (!$theme) {
        return "";    
    }
    return $theme->getName();
}
