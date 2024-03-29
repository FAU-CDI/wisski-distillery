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



/** get_footer_region returns the region that implements the footer */
function get_footer_region(): string {
  $footer_block_map = [
    "teriyaki" => "footer_bottom",
    "olivero" => "footer_bottom",
    "bartik" => "footer_fifth",
    "ffbartik" => "footer_fifth",
    "dxpr_theme" => "footer",
    "bootstrap" => "footer",
    "bootstrap4" => "footer",
    "bootstrap5" => "footer",
    "bootstrap_for_drupal" => "footer_sub_center",
    "bootstrap_for_drupal_subtheme" => "footer_sub_center",
    "roma_theme" => "footer_fifth",
    "gnm2018" => "footer",
    "oin_graphik" => "footer",
    "oin_paleo" => "footer",
    "oin_projekt" => "footer",
  ];

  $theme = \Drupal::service('theme.manager')->getActiveTheme();
  if (!$theme) {
    return "";
  }
  return $footer_block_map[$theme->getName()] ?? ""; // return the theme
}