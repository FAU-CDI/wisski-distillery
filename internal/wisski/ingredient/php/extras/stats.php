<?php

function export_statistics() {
    $service = 'wisski_statistics.statistics';
    if (empty(\Drupal::hasService($service))) return (object)[];
    $statistics = \Drupal::service($service);
    return $statistics->update();
}
