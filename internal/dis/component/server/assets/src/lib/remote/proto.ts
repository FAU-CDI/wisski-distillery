import { Mutex } from "async-mutex";

export type CallMessage = { name: string; params?: string[] | null }
export type ResultMessage = { success: boolean; message?: string; }
export type SignalMessage = { signal: string; }
function isResultMessage(value: any): value is ResultMessage {
    return typeof value === 'object' &&
        Object.prototype.hasOwnProperty.call(value, 'success') && typeof value['success'] === 'boolean' &&
        (!Object.prototype.hasOwnProperty.call(value, 'message') || (value['message'] === undefined || typeof value['message'] === 'string'));
}

/**
 * Opens a WebSocket connection and calls a server action
 * @param endpoint Endpoint to call
 * @param call Function to call
 * @param onOpen callback for once the connection is opened. The send function can be used to send additional text to the server.
 * @param onText called when the connection receives some text
 * @returns a promise that is resolved once the conneciton is closed. Rejected if the connection errors.
 */
export default async function callServerAction(
    endpoint: string,
    call: CallMessage,
    onOpen: (send: (text: string) => void, cancel: () => void) => void,
    onText: (text: string) => void,
): Promise<ResultMessage> {
    return new Promise((rs, rj) => {
        const mutex = new Mutex();

        const socket = new WebSocket(endpoint);

        let result: ResultMessage;
        socket.onmessage = (msg) => {
            mutex.runExclusive(async () => {
                if (typeof msg.data === 'string') {
                    onText(msg.data)
                    return
                }

                if (msg.data instanceof Blob) {
                    const object = JSON.parse(await msg.data.text())
                    if (isResultMessage(object)) {
                        result = {'success': object['success'], 'message': object['message']}
                        return
                    }
                }

                console.warn('Unknown message', msg)
            })
        }
        socket.onclose = () => {
            mutex.runExclusive(() => rs(result));
        };
        socket.onerror = rj; // if an error occurs, close the socket

        socket.onopen = () => {
            const blob = new Blob([JSON.stringify(call)]);
            socket.send(blob)

            onOpen(
                (text: string) => {
                    if (typeof text !== 'string') {
                        console.warn('Ignoring send() call with unknown type', text)
                        return
                    }
                    socket.send(text)
                },
                () => {
                    const blob = new Blob([JSON.stringify({'signal': 'cancel'})]);
                    socket.send(blob)
                }
            )
        }
    })
}