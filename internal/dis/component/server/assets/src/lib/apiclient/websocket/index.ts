/** @file implements the websocket protocol used by the distillery */

import WebSocket from 'isomorphic-ws'
import { Buffer } from 'buffer'

/** Call represents a specific WebSocket call */
export default class Call {
  constructor (remote: Remote, spec: CallSpec) {
    this.remote = remote
    this.call = spec
  }

  public readonly remote: Readonly<Remote>
  public readonly call: Readonly<CallSpec>

  /** called right before sending the request */
  public beforeCall?: (this: Call) => void

  /** called right after the socket is closed */
  public afterCall?: (this: Call, result: Result) => void

  /** called when an error occurs before rejecting the promise */
  public onError?: (this: Call, error: any) => void

  /** called when a log line is received */
  public onLogLine?: (this: Call, line: string) => void

  /** connect checks if the connect method was called */
  private connected: boolean = false

  /** holds the websocket when the connection is alive */
  private ws: WebSocket | null = null

  /**
   * Connect to the specified remote endpoint and perform the action
   * @param remote Remote to connect to
   */
  async connect (): Promise<Result> {
    // ensure that connect is only run once.
    if (this.connected) {
      throw new Error('connect() may only be called once')
    }
    this.connected = true

    // and do the connection!
    return await new Promise((resolve, reject) => {
      // create the websocket
      const ws = new WebSocket(this.remote.url, typeof this.remote.token === 'string' ? { headers: { Authorization: 'Bearer ' + this.remote.token } } : undefined)
      this.ws = ws // make it available to other things

      // result is a promise, because some APIs in the browser are async
      let result = Promise.resolve(JSON.stringify({ success: false, message: 'Unknown error' }))

      ws.onopen = () => {
        if (this.beforeCall != null) {
          this.beforeCall()
        }
        ws.send(Buffer.from(JSON.stringify(this.call), 'utf8'))
      }

      ws.onmessage = ({ data }) => {
        // if this is a string it is a log line
        if (typeof data === 'string') {
          if (this.onLogLine != null) {
            this.onLogLine(data)
          }
          return
        }

        // decode the message
        if (data instanceof Blob) {
          result = data.text()
        } else {
          const decoder = new TextDecoder()
          result = Promise.resolve(decoder.decode(data as ArrayBuffer))
        }
      }

      ws.onerror = (err) => {
        this.close()

        // call the handler and reject
        if (this.onError != null) {
          this.onError(err)
        }
        reject(err)
      }

      ws.onclose = () => {
        this.close()

        // decode the result
        result
          .then(t => JSON.parse(t))
          .then((res) => {
            if (this.afterCall != null) {
              this.afterCall(res)
            }
            resolve(res)
          })
          .catch((e) => console.error(e))
      }
    })
  }

  /** sendText sends some text to the server requests cancellation of an ongoing operation */
  sendText (text: string): void {
    const ws = this.ws
    if (ws == null) {
      throw new Error('websocket not connected')
    }

    ws.send(text)
  }

  /** cancel requests cancellation of an ongoing operation */
  cancel (): void {
    const ws = this.ws
    if (ws == null) {
      throw new Error('websocket not connected')
    }

    ws.send(Buffer.from(JSON.stringify({ signal: 'cancel' }), 'utf8'))
  }

  /** close closes this websocket connection */
  private close (): void {
    const ws = this.ws
    if (ws == null) {
      throw new Error('websocket not connected')
    }

    ws.close()
    this.ws = null
  }
}

/** specifies a remote endpoint */
export interface Remote {
  url: string // the remote websocket url to talk to
  token?: string // optional token
}

/** CallSpec represents the specification for a call */
export interface CallSpec {
  call: string
  params: string[]
}

/** the result of a websocket call */
export interface Result {
  success: boolean
  message: string
}
