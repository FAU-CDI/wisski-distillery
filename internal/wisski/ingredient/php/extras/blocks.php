<?php

/**
 * Creates a basic block and optionally places it into the right region
 */
function create_basic_block(string $info, string $html, string $region, string $block_id) {

  // create a custom block
  $block_content = \Drupal\block_content\Entity\BlockContent::create([
    'info' => $info,
    'type' => 'basic',
    'body' => [
      'value' => $html,
      'format' => 'full_html',
    ],
  ]);
  $block_content->save();

  if ($region === "") {
    return;
  }

  // get plugin and theme id
  $plugin = 'block_content:' . $block_content->uuid();
  $theme = \Drupal::theme()->getActiveTheme()->getName();

  $block = \Drupal\block\Entity\Block::create([
    'plugin' => $plugin,
    'id' => $block_id,
    'region' => $region,
    'status' => TRUE,
    'theme' => $theme,
    'weight' => 0,
  ]);
  $block->save();
}