export default function connectSocket(onOpen: (socket: WebSocket) => void, onData: (data: any) => void): Promise<CloseEvent> {
    return new Promise((rs, rj) => {
        const socket = new WebSocket(location.href.replace('http', 'ws'));

        socket.onclose = rs;
        socket.onerror = rj;

        socket.onmessage = (ev) => onData(ev.data)
        socket.onopen = () => onOpen(socket);
    });
}