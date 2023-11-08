/** @file implements the websocket protocol used by the distillery */

import WebSocket from "ws";

/** A call to the websocket endpoint */
export interface WebSocketCall {
  call: string;
  params: string[];
}

/** the result of a websocket call */
export interface WebSocketResult {
  success: boolean,
  message: string,
}

/** optional hooks to call when something happens */
export interface Hooks {
  beforeCall: (call: WebSocketCall) => void;   // called right before sending the request
  afterCall: (call: WebSocketCall, result: WebSocketResult) => void;   // called when the socket is closed
  onError: (call: WebSocketCall, error: any) => void;   // called when an error occurs before rejecting the promise
  onLogLine: (call: WebSocketCall, line: string) => void;   // called when a log line is received
}

/** specifies a remote endpoint */
export interface Remote {
  url: string; // the remote websocket url to talk to
  token?: string; // optional token
}

/** run a websocket remote call */
export default async function Call(remote: Remote, call: WebSocketCall, hooks?: Partial<Hooks>): Promise<WebSocketResult> {
  return new Promise((resolve, reject) => {
    let options = { headers: {} };
    if (remote.token) {
      options.headers = { 'Authorization': 'Bearer ' + remote.token };
    }
    const ws = new WebSocket(remote.url, options);

    let result = {'success': false, 'message': 'Unknown error'};
    ws.on('error', (err) => {
      if (hooks && hooks.onError) {
        hooks.onError(call, err);
      }
      reject(err)
    });
    ws.on('open', () => {
      if (hooks && hooks.beforeCall) {
        hooks.beforeCall(call);
      }
      ws.send(Buffer.from(JSON.stringify(call), 'utf8'));
    });

    ws.on('message', async (msg, isBinary) => {
      if (!isBinary) {
        if (hooks && hooks.onLogLine) {
          hooks.onLogLine(call, msg.toString());
        }
        return;
      }
      result = JSON.parse(msg.toString());
    });

    ws.on('close', () => {
      if (hooks && hooks.afterCall) {
        hooks.afterCall(call, result);
      }
      resolve(result);
    });
  });
}
