import { createModal } from '~/src/lib/remote'

/**
 * Flags to provision a new system.
 * Should mirror "provision".Flags.
 */
interface ProvisionFlags {
  Slug: string
  Flavor?: string
  System: System
}

interface System {
  PHP: string
  OpCacheDevelopment: boolean
  ContentSecurityPolicy: string
}

/** Rebuild the specified instance */
export async function Rebuild (slug: string, system: System): Promise<string> {
  return await new Promise((resolve, reject) => {
    createModal('rebuild', [slug, JSON.stringify(system)], {
      bufferSize: 0,
      onClose: (success: boolean, message?: string) => {
        if (!success) {
          reject(new Error(message ?? 'unspecified error'))
          return
        }

        resolve(slug)
      }
    })
  })
}

/** Provision provisions a new instance */
export async function Provision (flags: ProvisionFlags): Promise<string> {
  // open a modal to provision a new instance
  return await new Promise((resolve, reject) => {
    createModal('provision', [JSON.stringify(flags)], {
      bufferSize: 0,
      onClose: (success: boolean, message?: string) => {
        if (!success) {
          reject(new Error(message ?? 'unspecified error'))
          return
        }

        resolve(flags.Slug)
      }
    })
  })
}
