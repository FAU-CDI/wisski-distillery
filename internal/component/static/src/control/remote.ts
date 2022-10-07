import connectSocket from '../socket/socket';

const elements = document.getElementsByClassName('remote-action')
Array.from(elements).forEach((element) => {
    const action = element.getAttribute('data-action') as string;
    const param = element.getAttribute('data-param') as string | undefined;
    const target = document.querySelector(element.getAttribute('data-target')!) as HTMLElement;
    const bufferSize = (function() {
        const number = parseInt(element.getAttribute('data-buffer') ?? "", 10) ?? 0;
        return (isFinite(number) && number > 0) ? number : 0; 
    })()

    let running = false
    element.addEventListener('click', function(ev) {
        ev.preventDefault();
        
        // already running
        if (running) return

        running = true
        element.setAttribute('disabled', 'disabled');
        const close = function() {
            element.removeAttribute('disabled');
            running = false;
        }

        target.innerText = "";

        const buffer: Array<string> = [];
        const println = function(line: string) {
            if(bufferSize === 0) {
                target.innerText += line + "\n";
                return;
            }
            
            buffer.push(line);
            if(buffer.length > bufferSize) {
                buffer.splice(0, buffer.length - bufferSize)
            }
            target.innerText = buffer.join("\n");
        }
        
        println("Connecting ...")

        // connect to the socket and send the action
        connectSocket((socket) => {
            println("Connected")
            socket.send(action);
            if (typeof param === 'string') {
                socket.send(param);
            }
        }, (data) => {
            println(data);
        }).then(() => {
            println("Connection closed.\n")
            close();
        }).catch(() => {
            println("Connection errored.\n")
            close();
        });
    });
})