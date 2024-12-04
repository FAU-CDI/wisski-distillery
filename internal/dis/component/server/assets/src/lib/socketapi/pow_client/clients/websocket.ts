/** @file implements the websocket protocol used by the distillery */

import WebSocket from 'modern-isomorphic-ws'
import { Buffer } from 'buffer'
import { type Session, type CallSpec, type Remote, type Result, WaitResult, isResult } from '../common/types'
import { Lazy } from '../common/once'
import { errAlreadyConnected, errNotConnected } from '../common/errors'

const EXIT_STATUS_NORMAL_CLOSE = 1000
const PROTOCOL = 'pow-1'

/**
 * A process-over-websocket session via the websocket-based protocol
 */
export default class WebsocketSession implements Session {
  constructor (public readonly remote: Remote, public readonly call: CallSpec) {
  }

  /** called right before sending the request */
  public beforeCall?: (this: WebsocketSession) => void

  /** called right after the socket is closed */
  public afterCall?: (this: WebsocketSession, result: Result) => void

  /** called when an error occurs before rejecting the promise */
  public onError?: (this: WebsocketSession, error: any) => void

  /** called when a log line is received */
  public onLogLine?: (this: WebsocketSession, line: string) => void

  /** holds the websocket when the connection is alive */
  private ws: WebSocket | null = null

  /** connect checks if the connect method was called */
  #connected: boolean = false

  private static decoder = new TextDecoder()

  async connect (): Promise<void> {
    // ensure that connect is only run once.
    if (this.#connected) {
      throw errAlreadyConnected
    }
    this.#connected = true

    let receivedResult: Result | null = null

    await new Promise<void>((_resolve, _reject) => {
      this.#result.Get(
        async () => await new Promise<Result>((resolve, reject) => {
          // create the websocket
          const ws = new WebSocket(this.remote.url, PROTOCOL, { headers: this.remote.headers })
          ws.binaryType = 'arraybuffer'
          this.ws = ws // make it available to other thing

          ws.onopen = () => {
            this.#closeStateHack()

            if (typeof this.beforeCall === 'function') {
              this.beforeCall()
            }

            this.#send(Buffer.from(JSON.stringify(this.call), 'utf8')).then(_resolve).catch(_reject)
          }

          ws.onmessage = ({ data, ...rest }: { data: unknown }) => {
            // ignore non-strings for now
            if (typeof data !== 'string') {
              if (!(data instanceof ArrayBuffer)) {
                console.error('websocket implementation did not return a buffer or blob')
                return
              }

              try {
                const raw = JSON.parse(WebsocketSession.decoder.decode(data))
                if (!isResult(raw, false)) {
                  throw new Error('result message is not an object')
                }
                receivedResult = raw
              } catch(err: unknown) {
                  console.error('failed to decode result message', err, data)
              }
              return
            }

            if (this.onLogLine != null) {
              this.onLogLine(data)
            }
          }

          ws.onerror = (err: unknown) => {
            this.close()
            
            // reject both promised
            _reject(err)
            reject(err)
          }

          ws.onclose = (event: { code: number, reason: string, wasClean: boolean }) => {
            // if the connection was not clean, reject with an error
            if (!event.wasClean) {
              reject(new Error('unclean exit with code ' + event.code.toString() + ' ' + event.reason))
              return
            }

            // normal close => process succeeded
            if (event.code !== EXIT_STATUS_NORMAL_CLOSE) {
              reject(new Error('unexpected exit with code ' + event.code.toString() + ' ' + event.reason))
              return
            }
            if (event.reason !== '') {
              reject(new Error('websocket closed with non-empty reason field'))
              return
            }

            if (receivedResult === null) {
              reject(new Error('did not receive a result'))
              return
            }

            resolve(receivedResult)
            this.close()
          }
        })
      )
    })
  }

  #result = new Lazy<Result>()
  async wait (): Promise<WaitResult> {
    const result = await this.#result.Get(() => { throw new Error('never reached') })
    return { result }
  }

  /**
   * Sometimes for unknown reasons the websocket gets stuck in CLOSING state.
   *
   * This code triggers code to manually unstick the server
   */
  #closeStateHack (): void {
    const STATE_POLL_INTERVAL = 100 // how often to poll the state
    const CLOSE_TIMEOUT = 500 // how long to wait for the close to finish on it's own

    const poller = setInterval(() => {
      // if we have an open or connecting websocket keep going
      const ws = this.ws
      if (ws !== null && (ws.readyState === ws.OPEN || ws.readyState === ws.CONNECTING)) {
        return
      }

      // clear the interval and only continue if in CLOSING state
      clearInterval(poller)
      if (ws === null || ws.readyState !== ws.CLOSING) {
        return
      }

      setTimeout(() => {
        if (ws.readyState === ws.CLOSING) {
          console.warn('websocket client misbehaved: still in closing state')
          ws.terminate()
        }
      }, CLOSE_TIMEOUT)
    }, STATE_POLL_INTERVAL)
  }

  #inputClosed = false
  /** sendText sends some text to the server requests cancellation of an ongoing operation */
  async sendText (text: string): Promise<void> {
    if (this.#inputClosed) return
    await this.#send(text)
  }

  /** cancel requests cancellation of an ongoing operation */
  async cancel (): Promise<void> {
    await this.#send(Buffer.from(JSON.stringify({ signal: 'cancel' }), 'utf8'))
  }

  /**
   * closeInput closes the input from the client
   * Any further text received on the server side will be ignored.
   */
  async closeInput (): Promise<void> {
    if (this.#inputClosed) {
      return
    }
    this.#inputClosed = true
    await this.#send(Buffer.from(JSON.stringify({ signal: 'close' }), 'utf8'))
  }

  static readonly #useSyncronousSend = (typeof window !== 'undefined')
  async #send (data: string | Buffer): Promise<void> {
   if (WebsocketSession.#useSyncronousSend) {
    return await this.#sendSync(data)
   }
   return this.#sendAsync(data)
  }

  async #sendAsync(data: string | Buffer): Promise<void> {
    const ws = this.ws
    if (ws == null) {
      throw errNotConnected
    }

    await new Promise<void>((resolve, reject) => {
      ws.send(data, {}, (err: Error | undefined) => {
        if (typeof err !== 'undefined' && err !== null) {
          reject(err)
          return
        }
        resolve()
      })
    })
  }

  async #sendSync(data: string | Buffer): Promise<void> {
    const ws = this.ws
    if (ws == null) {
      throw errNotConnected
    }

    await new Promise<void>(resolve => {
      ws.send(data)
      resolve()
    })
  }

  /** close closes this websocket connection */
  private close (): void {
    const ws = this.ws
    if (ws == null) {
      throw errNotConnected
    }

    ws.close()
    this.ws = null
  }
}
