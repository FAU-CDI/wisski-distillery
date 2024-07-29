import { CallSpec, WebsocketSession } from "../socketapi/pow_client";

/** LocalCall is like Call, but uses the current page */
export default class LocalSession extends WebsocketSession {
  constructor (spec: CallSpec) {
    super({ url: location.protocol.replace('http', 'ws') + '//' + location.host + '/api/v1/pow/' }, spec)
  }
}
