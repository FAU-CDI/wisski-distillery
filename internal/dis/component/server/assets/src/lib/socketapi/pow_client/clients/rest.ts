import { type Session, type CallSpec, type Remote, WaitResult, Status, isStatus } from '../common/types'
import axios from 'axios'
import { sleep } from '../common/utils'
import { Lazy } from '../common/once'
import { errAlreadyConnected, errNotConnected, errWaitExceeded } from '../common/errors'

/**
 * A process-over-websocket session via the REST-based protocol
 */
export default class RestSession implements Session {
  constructor (public readonly remote: Remote, public readonly call: CallSpec) {
  }

  #connected = false
  #id: string | null = null

  async connect (): Promise<void> {
    if (this.#connected) {
      throw errAlreadyConnected
    }
    this.#connected = true

    const id = await this.#rest('/new', this.call)
    if (typeof id !== 'string') {
      throw new Error('did not receive an id back')
    }
    this.#id = id
  }

  readonly #result = new Lazy<WaitResult>()
  async wait (options?: { pollInterval?: number, maxWait?: number }): Promise<WaitResult> {
    return await this.#result.Get(async (): Promise<WaitResult> => {
      const pollInterval = options?.pollInterval ?? 500
      const maxWait = options?.maxWait ?? 60 * 60 * 1000 // defaults to 1 hour

      const start = performance.now()
      let status: Status = await this.status()
      while (status.result.status === 'pending') {
        if (performance.now() - start > maxWait) {
          throw errWaitExceeded
        }
        await sleep(pollInterval)
        status = await this.status()
      }
      return status as WaitResult 
    })
  }

  /**
   * Fetches the current connection state from the server.
   * If the connection is not connected, throws {@link errNotConnected}.
   */
  async status (): Promise<Status> {
    if (typeof this.#id !== 'string') throw errNotConnected

    const status = await this.#rest(`/status/${this.#id}`)
    if (!isStatus(status)) {
      return {'result': {'status': 'rejected', 'reason': 'invalid status returned'}}
    }
    return status
  }

  #inputClosed = false
  async sendText (text: string): Promise<void> {
    if (typeof this.#id !== 'string') throw errNotConnected
    if (this.#inputClosed) return
    return await this.#rest(`/input/${this.#id}`, text + '\n')
  }

  async cancel (): Promise<void> {
    if (typeof this.#id !== 'string') throw errNotConnected
    return await this.#rest(`/cancel/${this.#id}`, null)
  }

  async closeInput (): Promise<void> {
    if (typeof this.#id !== 'string') throw errNotConnected
    if (this.#inputClosed) return
    this.#inputClosed = true
    return await this.#rest(`/closeInput/${this.#id}`, null)
  }

  /** sends a get request if data is not provided, and a post request otherwise */
  async #rest (path: string, data?: any): Promise<any> {
    const url = this.#buildURL(path)

    const { headers } = this.remote
    const config = (typeof headers !== 'undefined') ? { headers } : undefined

    const res = await ((typeof data !== 'undefined') ? axios.post(url, data, config) : axios.get(url, config))
    if (res.status !== 200) {
      throw new Error('received invalid status code')
    }
    return res.data
  }

  /** builds the URL for the client to connect to the specified path */
  #buildURL (path: string): string {
    const { url: base } = this.remote

    const baseNoSlash = base.endsWith('/') ? base.substring(0, base.length - 1) : base
    const pathNoSlash = path.startsWith('/') ? path.substring(1) : path

    return baseNoSlash + '/' + pathNoSlash
  }
}
