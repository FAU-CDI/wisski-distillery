export { default as WebsocketSession } from './clients/websocket'
export { default as RestSession } from './clients/rest'

export { errAlreadyConnected, errNotConnected, errWaitExceeded } from './common/errors'
export type { Remote, CallSpec, Result, Session } from './common/types'
