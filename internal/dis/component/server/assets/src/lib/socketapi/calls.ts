/** @file provides a list of websocket calls supported by the backend */

import { CallSpec } from "./pow_client"

/** Backup backups everything */
export function Backup (): CallSpec {
  return {
    call: 'backup',
    params: []
  }
}

interface ProvisionParams {
  Slug: string
  Flavor?: 'Drupal 10' | 'Drupal 9'
  IIPServer?: string
  System: SystemParams
}

interface SystemParams {
  PHP: 'Default (8.1)' | '8.0' | '8.1' | '8.2' | '8.3'
  OpCacheDevelopment: boolean
  ContentSecurityPolicy: string
}

/** Provision provisions a new instance */
export function Provision (params: ProvisionParams): CallSpec {
  return {
    call: 'provision',
    params: [
      JSON.stringify(params)
    ]
  }
}

/** Snapshot makes a snapshot of an instance */
export function Snapshot (Slug: string): CallSpec {
  return {
    call: 'snapshot',
    params: [Slug]
  }
}

/** Rebuild rebuilds an instance */
export function Rebuild (Slug: string, params: SystemParams): CallSpec {
  return {
    call: 'rebuild',
    params: [
      Slug,
      JSON.stringify(params)
    ]
  }
}

/** Update updates a specific instance */
export function Update (Slug: string): CallSpec {
  return {
    call: 'update',
    params: [Slug]
  }
}

/** Start starts a specific instance */
export function Start (Slug: string): CallSpec {
  return {
    call: 'start',
    params: [Slug]
  }
}

/** Stop stops a specific instance */
export function Stop (Slug: string): CallSpec {
  return {
    call: 'stop',
    params: [Slug]
  }
}

/** Purge purges a specific instance */
export function Purge (Slug: string): CallSpec {
  return {
    call: 'purge',
    params: [Slug]
  }
}
