import { default as Call, CallSpec } from '../apiclient/websocket';

/** LocalCall is like Call, but uses the current page */
export default class LocalCall extends Call {
    constructor(spec: CallSpec) {
        super({ url: location.protocol.replace('http', 'ws') + '//' + location.host + '/api/v1/ws'}, spec);
    }
}