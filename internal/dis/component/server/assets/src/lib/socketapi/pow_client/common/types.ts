import { type errAlreadyConnected, type errNotConnected, type errWaitExceeded } from './errors' // eslint-disable-line @typescript-eslint/no-unused-vars

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
  success: true
  data: unknown
  buffer?: string
}

interface ResultFailure {
  success: false
  data: string // error message (if any)
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
  wait: (options?: { pollInterval?: number, maxWait?: number }) => Promise<Result>

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
