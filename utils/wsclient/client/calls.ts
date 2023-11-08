import type { WebSocketCall } from ".";

/** Backup backups everything */
export function Backup(): WebSocketCall {
    return {
        'call': 'backup',
        'params': [],
    } 
}

type ProvisionParams = {
    Slug: string;
    Flavor?: "Drupal 10" | "Drupal 9",
    System: SystemParams
}

type SystemParams = {
    PHP: "Default (8.1)" | "8.0" | "8.1" | "8.2",
    OpCacheDevelopment: boolean,
    ContentSecurityPolicy: string,
}

/** Provision provisions a new instance */
export function Provision(params: ProvisionParams): WebSocketCall {
    return {
        'call': 'provision',
        'params': [
            JSON.stringify(params)
        ],
    } 
}

/** Snapshot makes a snapshot of an instance */
export function Snapshot(Slug: string): WebSocketCall {
    return {
        'call': 'snapshot',
        'params': [Slug],
    }
}

/** Rebuild rebuilds an instance */
export function Rebuild(Slug: string, params: SystemParams): WebSocketCall {
    return {
        'call': 'rebuild',
        'params': [
            Slug,
            JSON.stringify(params)
        ],
    } 
}

/** Update updates a specific instance */
export function Update(Slug: string): WebSocketCall {
    return {
        'call': 'update',
        'params': [Slug],
    }
}


/** Start starts a specific instance */
export function Start(Slug: string): WebSocketCall {
    return {
        'call': 'start',
        'params': [Slug],
    }
}

/** Stop stops a specific instance */
export function Stop(Slug: string): WebSocketCall {
    return {
        'call': 'stop',
        'params': [Slug],
    }
}

/** Purge purges a specific instance */
export function Purge(Slug: string): WebSocketCall {
    return {
        'call': 'purge',
        'params': [Slug],
    }
}
