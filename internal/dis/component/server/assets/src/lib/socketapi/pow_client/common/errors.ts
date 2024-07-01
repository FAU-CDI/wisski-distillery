/** Indicates that the connection has already been established */
export const errAlreadyConnected = new Error('connection already established')

/** Indicates that the connection is not connected */
export const errNotConnected = new Error('not connected')

/** Indicates that the wait call ran out of time */
export const errWaitExceeded = new Error('wait time limit exceeded')
