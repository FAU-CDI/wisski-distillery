package httpx

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tkw1536/pkglib/lazy"
)

// WebSocket implements serving a WebSocket
type WebSocket struct {
	Context context.Context // context which closes all connections
	Limits  WebSocketLimits // limits for websocket operations

	Handler  func(ws WebSocketConnection)
	Fallback http.Handler

	pool     lazy.Lazy[*sync.Pool] // pool holds *WebSocketConn objects
	upgrader websocket.Upgrader    // upgrades upgrades connections
}

type WebSocketLimits struct {
	WriteWait      time.Duration // maximum time to wait for writing
	PongWait       time.Duration // time to wait for pong responses
	PingInterval   time.Duration // interval to send pings to the client
	MaxMessageSize int64         // maximal message size in bytes
}

func (limits *WebSocketLimits) SetDefaults() {
	if limits.WriteWait == 0 {
		limits.WriteWait = 10 * time.Second
	}
	if limits.PongWait == 0 {
		limits.PongWait = time.Minute
	}
	if limits.PingInterval <= 0 {
		limits.PingInterval = (limits.PongWait * 9) / 10
	}
	if limits.MaxMessageSize <= 0 {
		limits.MaxMessageSize = 2048
	}
}

// makePoolSocket creates a new socket and makes sure that the pool is initialized
func (h *WebSocket) makePoolSocket() *webSocketConn {
	return h.pool.Get(func() *sync.Pool {
		return &sync.Pool{
			New: func() any { return new(webSocketConn) },
		}
	}).Get().(*webSocketConn)
}

func (h *WebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// if the user did not request a websocket, go to the fallbacjk handler
	if !websocket.IsWebSocketUpgrade(r) {
		h.serveFallback(w, r)
		return
	}

	// else deal with the websocket!
	h.serveWebsocket(w, r)
}

func (h *WebSocket) serveFallback(w http.ResponseWriter, r *http.Request) {
	if h.Fallback == nil {
		http.NotFound(w, r)
		return
	}

	h.Fallback.ServeHTTP(w, r)
}

func (h *WebSocket) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	// upgrade the connection or bail out!
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// get a new socket from the pool
	socket := h.makePoolSocket()
	socket.Serve(h.Context, h.Limits, conn, h.Handler)

	// return a reset socket to the pool
	socket.reset()
	h.pool.Get(nil).Put(socket)
}

// WebSocketConnection represents a connected WebSocket.
type WebSocketConnection interface {
	// Context returns a context that is closed once the connection is terminated.
	Context() context.Context

	// Read returns a channel that receives message.
	// The channel is closed if no more messags are available (for instance because the server closed).
	Read() <-chan WebSocketMessage

	// Write queues the provided message for sending.
	// The returned channel is closed once the message has been sent.
	Write(WebSocketMessage) <-chan struct{}

	// WriteText is a convenience method to send a TextMessage.
	// The returned channel is closed once the message has been sent.
	WriteText(text string) <-chan struct{}

	// Close closes the underlying connection
	Close()
}

// WebSocketMessage represents a connected Websocket
type WebSocketMessage struct {
	Type  int
	Bytes []byte
}

type outWebSocketMessage struct {
	WebSocketMessage
	done chan<- struct{} // done should be closed when finished
}

// webSocketConn implements [WebSocketConnection]
type webSocketConn struct {
	conn   *websocket.Conn // underlying connection
	limits WebSocketLimits

	context context.Context // context to cancel the connection
	cancel  context.CancelFunc

	wg sync.WaitGroup // blocks all the ongoing tasks

	// incoming and outgoing tasks
	incoming chan WebSocketMessage
	outgoing chan outWebSocketMessage
}

// Serve serves the provided connection
func (h *webSocketConn) Serve(ctx context.Context, limits WebSocketLimits, conn *websocket.Conn, handler func(ws WebSocketConnection)) {
	// use the connection!
	h.conn = conn

	// setup limits
	h.limits = limits
	h.limits.SetDefaults()

	// create a context for the connection
	if ctx == nil {
		ctx = context.Background()
	}
	h.context, h.cancel = context.WithCancel(ctx)

	// start receiving and sending messages
	h.wg.Add(2)
	h.sendMessages()
	h.recvMessages()

	// wait for the context to be cancelled, then close the connection
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		<-h.context.Done()
		h.conn.Close()
	}()

	// start the application logic
	h.wg.Add(1)
	go h.handle(handler)

	// wait for closing operations
	h.wg.Wait()
}

func (h *webSocketConn) handle(handler func(ws WebSocketConnection)) {
	defer func() {
		h.wg.Done()
		h.cancel()
	}()

	handler(h)
}

func (h *webSocketConn) sendMessages() {
	h.outgoing = make(chan outWebSocketMessage)

	go func() {
		// close connection when done!
		defer func() {
			h.wg.Done()
			h.cancel()
		}()

		// setup a timer for pings!
		ticker := time.NewTicker(h.limits.PingInterval)
		defer ticker.Stop()

		for {
			select {
			// everything is done!
			case <-h.context.Done():
				return

			// send outgoing messages
			case message := <-h.outgoing:
				(func() {
					defer close(message.done)

					err := h.writeRaw(message.Type, message.Bytes)
					if err != nil {
						return
					}
					message.done <- struct{}{}
				})()
			// send a ping message
			case <-ticker.C:
				if err := h.writeRaw(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			}
		}
	}()

}

// writeRaw writes to the underlying socket
func (h *webSocketConn) writeRaw(messageType int, data []byte) error {
	h.conn.SetWriteDeadline(time.Now().Add(h.limits.WriteWait))
	return h.conn.WriteMessage(messageType, data)
}

// Write writes a message to the websocket connection.
func (sh *webSocketConn) Write(message WebSocketMessage) <-chan struct{} {
	callback := make(chan struct{}, 1)
	go func() {
		select {
		// write an outgoing message
		case sh.outgoing <- outWebSocketMessage{
			WebSocketMessage: message,
			done:             callback,
		}:
		// context
		case <-sh.context.Done():
			close(callback)
		}
	}()
	return callback
}

func (sh *webSocketConn) WriteText(text string) <-chan struct{} {
	return sh.Write(WebSocketMessage{
		Type:  websocket.TextMessage,
		Bytes: []byte(text),
	})
}

func (h *webSocketConn) recvMessages() {
	h.incoming = make(chan WebSocketMessage)

	// set a read handler
	h.conn.SetReadLimit(h.limits.MaxMessageSize)

	// configure a pong handler
	h.conn.SetReadDeadline(time.Now().Add(h.limits.PongWait))
	h.conn.SetPongHandler(func(string) error { h.conn.SetReadDeadline(time.Now().Add(h.limits.PongWait)); return nil })

	// handle incoming messages
	go func() {
		// close connection when done!
		defer func() {
			h.wg.Done()
			h.cancel()
		}()

		for {
			messageType, messageBytes, err := h.conn.ReadMessage()
			if err != nil {
				return
			}

			// try to send a message to the incoming message channel
			select {
			case h.incoming <- WebSocketMessage{
				Type:  messageType,
				Bytes: messageBytes,
			}:
			case <-h.context.Done():
				return
			}
		}
	}()
}

// Read returns a channel that receives incoming messages.
// The channel is close once no more messages are available, or the context is canceled.
func (h *webSocketConn) Read() <-chan WebSocketMessage {
	return h.incoming
}

// Context returns a context that is closed once this connection is closed.
func (h *webSocketConn) Context() context.Context {
	return h.context
}

func (h *webSocketConn) Close() {
	h.cancel()
}

// reset resets this websocket
func (h *webSocketConn) reset() {
	h.limits = WebSocketLimits{}
	h.conn = nil
	h.incoming = nil
	h.outgoing = nil
	h.context, h.cancel = nil, nil
}
