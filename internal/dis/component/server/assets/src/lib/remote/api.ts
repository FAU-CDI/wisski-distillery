import { createModal } from "~/src/lib/remote"

/**
 * Flags to provision a new system.
 * Should mirror "provision".Flags.
 */
interface ProvisionFlags {
    Slug: string
    System: System
}

interface System {
    PHP: string;
    OpCacheDevelopment: boolean
}

/** Rebuild the specified instance */
export async function Rebuild(slug: string, system: System): Promise<string> {
    return new Promise((rs, rj) => {
        createModal("rebuild", [slug, JSON.stringify(system)], {
            bufferSize: 0,
            onClose: (success: boolean, message?: string) => {
                if (!success) {
                    rj(new Error(message ?? "unspecified error"))
                    return;
                }
                
                rs(slug);
            },
        })
    });
}

/** Provision provisions a new instance */
export async function Provision(flags: ProvisionFlags): Promise<string> {
    // open a modal to provision a new instance
    return new Promise((rs, rj) => {
        createModal("provision", [JSON.stringify(flags)], {
            bufferSize: 0,
            onClose: (success: boolean, message?: string) => {
                if (!success) {
                    rj(new Error(message ?? "unspecified error"))
                    return;
                }
                
                rs(flags.Slug);
            },
        })
    });
}