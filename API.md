# API Documentation

The distillery comes with an API served under `/api/`.
There are two disinct types of APIs:
- a json-based GET/POST API for quick information retrieval
- a websocket-based API for interactive actions
The complete API is still a work in progress, and will be polished and properly implemented at a later point.
The API is currently disabled by default, and needs to be enabled in `distillery.yaml`. 

## JSON-based API

These routes use a simple REST-based protocol for GET and POST requests.
Requests are sent using either GET (for immutable request) or POST or DELETE (for mutable requests).
All requests should respond nearly instantly, returning JSON-encoded data to the client.
Typically each request takes only a second to execute.

NOTE: These routes will be documented using a Swagger / OpenAPI definition in the future. 
All routes can be found under `/api/v1/http/`

- `/api/v1/auth`: Returns api session information
- `/api/v1/news`: Returns JSON containing all news items
- `/api/v1/instances/directory`: Returns a (publically visible) list of systems 
- `/api/v1/resolve?uri=...`: Resolve a URI


## Interactive Websocket API

Some API calls require interactivity or provide streaming content to clients.
An example of such an action is creating a new instance.
The protocol is based on [Websockets](https://websockets.spec.whatwg.org/).
The API is reachable under `/api/v1/ws/`

The websocket API uses two kinds of frames:
- a `binary frame` to send JSON-encoded control data for the API; and
- a `text frame` to send output from or to (client -> server) the running API call.

Binary frames are utf8-encoded json serialization of objects.
These are used for bi-directional control data.

The client may send two kinds of binary frames:
- A `CallMessage`: ```json {"call": "some-name", "params": ["some", "string", "arguments"]}```. Different calls and their arguments are documented below. These must be sent exactly once when initializing the connection to indicate to the API which process should be started. 
- A `SignalMessage` of type "cancel": ```json {"signal": "cancel"}``` This may be sent at any time after the initial CallMessage to cancel the ongoing event. The ongoing event will be cancelled gracefully and perform cleanup.

If the client closes the connection (or loses contact), it is interpreted the same as a cancel signal.

The server may send two kinds of binary frames:
- A `ResultMessage` of type "success": ```json {"success": true}``` indicates that the process has finished successfully.
- A `ResultMessage` of type "failure": ```json {"success": false, "message": "some description of the failure"} indicates that the underlying process has failed with the given reason.

The server will always send exactly one of these frames before closing the connection.

Text frames contain input and output from the underlying process.
They may be sent at any point after the client sends the initial `CallMessage` and before the final `ResultMessage`.
Each frame contains data for a single line, not including the line terminators.
Each input is sent directly to the underlying process.

### Supported Websocket Calls

(to be documented)
