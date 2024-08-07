import { type errAlreadyConnected, type errNotConnected, type errWaitExceeded } from './errors'

/**
 * Specifies a remote endpoint for either protocol to connect to.
 */
export interface Remote {
  url: string // the remote websocket url to talk to
  headers?: Record<string, string>
}

/**
 * CallSpec is a single remote call
 */
export interface CallSpec {
  call: string
  params: string[]
}


/**
 * Result is the result of a websocket call
 */
export type Result = ResultSuccess | ResultFailure

interface ResultSuccess {
  status: 'fulfilled'
  value?: unknown
}

interface ResultFailure {
  status: 'rejected'
  reason?: string // error message (if any)
}

interface ResultPending {
  status: 'pending'
}

export interface WaitResult {
  result: Result
  buffer?: string
}

export interface Status {
  result: Result | ResultPending 
  buffer?: string
}

export function isStatus(value: unknown): value is Status {
  // must be object
  if (typeof value !== 'object' || value === null) return false

  // result must be a Result
  if (!('result' in value) || !isResult(value.result, true)) {
    return false
  }

  // buffer must be a string or undefined
  if ('buffer' in value) {
    return typeof value.buffer === 'string' || typeof value.buffer === 'undefined'
  }
  return true
}

export function isResult(value: unknown, allowPending: true): value is (Result | ResultPending)
export function isResult(value: unknown, allowPending: false): value is Result
export function isResult(value: unknown, allowPending: boolean): value is (Result | ResultPending) {
  // must be object
  if (typeof value !== 'object' || value === null) {
    return false
  }

  // status must exist 
  if (!('status' in value)) {
    return false
  }
  
  // status must be one of the allowed values
  switch(value['status']) {
    case 'fulfilled':
      return true
    case 'rejected':
      if ('reason' in value) {
        const reason = typeof value.reason
        return reason === 'string' || reason === 'undefined'
      }
      return true
    case 'pending':
      return allowPending
    default:
      return false  
  }
}

/** POWSession represents a session to connect to a remote */
export interface Session {
  readonly remote: Readonly<Remote>
  readonly call: Readonly<CallSpec>

  /**
   * Establishes a connection with the server and instructs
   * it to execute the given call.
   *
   * If the connection is already connected, throws {@link errAlreadyConnected}.
   */
  connect: () => Promise<void>

  /**
   * Wait waits for the session to complete and returns the result.
   * Implementations may chose to ignore pollInterval and maxWait parameters.
   *
   * If the connection is not connected, throws {@link errNotConnected}.
   * If the time limit is exceeded, throws {@link errWaitExceeded}.
   */
  wait: (options?: { pollInterval?: number, maxWait?: number }) => Promise<WaitResult>

  /**
   * Send text sends some text input to the connection.
   *
   * If the input is already closed, this is a no-op.
   * If the connection is not connected, throws {@link errNotConnected}.
   *
   * @returns
   */
  sendText: (text: string) => Promise<void>

  /**
   * Closes the input stream of the ongoing process.
   * Future calls to sendText return immediately.
   *
   * If the connection is not connected, throws {@link errNotConnected}.
   */
  closeInput: () => Promise<void>

  /**
   * Cancels an open connection.
   *
   * If the connection is not connected, throws {@link errNotConnected}.
   */
  cancel: () => Promise<void>
}
