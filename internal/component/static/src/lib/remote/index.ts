import "./index.css"
import connectSocket from './socket';

type Println = ((line: string, flush?: boolean) => void) & {
    paintedFrames: number;
    missedFrames: number;
}

/**
 * makeTextBuffer returns a println() function that efficiently writes text into target, and keeps at most size elements in the traceback.
 * scrollContainer is used to scroll on every painted update.
 */
function makeTextBuffer(target: HTMLElement, scrollContainer: HTMLElement, size: number): Println {
    let lastAnimationFrame: number | null = null; // last scheduled animation frame

    const buffer: Array<string> = []; // the internal buffer of lines
    const paint = () => {
        println.paintedFrames++
        target.innerText = buffer.join("\n")
        scrollContainer.scrollTop = scrollContainer.scrollHeight
        lastAnimationFrame = null
    }

    const println = (line: string, flush?: boolean) => {
        // add the line 
        buffer.push(line)
        if (size !== 0 && buffer.length > size) {
            buffer.splice(0, buffer.length - size)
        }

        // and update the browser in the next animation frame
        if (lastAnimationFrame !== null) {
            println.missedFrames++
            window.cancelAnimationFrame(lastAnimationFrame)
        }

        // force a repaint!
        if(flush) return paint(); 

        // schedule an animation frame
        lastAnimationFrame = window.requestAnimationFrame(paint);
    }
    println.paintedFrames = 0;
    println.missedFrames = 0;

    return println;
}

const elements = document.getElementsByClassName('remote-action')
Array.from(elements).forEach((element) => {
    const action = element.getAttribute('data-action') as string;
    const param = element.getAttribute('data-param') as string | undefined;
    const bufferSize = (function () {
        const number = parseInt(element.getAttribute('data-buffer') ?? "", 10) ?? 0;
        return (isFinite(number) && number > 0) ? number : 0;
    })()

    element.addEventListener('click', function (ev) {
        ev.preventDefault();

        // create a modal dialog and append it to the body
        const modal = document.createElement("div")
        modal.className = "modal-terminal"
        document.body.append(modal)

        // create a <pre> to write stuff into
        const target = document.createElement("pre")
        const println = makeTextBuffer(target, modal, bufferSize)
        modal.append(target)

        
        // create a button to eventually close everything
        const button = document.createElement("button")
        button.className = "pure-button pure-button-success"
        button.append("Close")
        button.addEventListener('click', function (event) {
            event.preventDefault();
            modal.parentNode?.removeChild(modal);
        })
        
        const onbeforeunload = window.onbeforeunload;
        window.onbeforeunload = () => "A remote session is in progress. Are you sure you want to leave?";

        // when closing, add a button to the modal!
        let didClose = false
        const close = function () {
            if (didClose) return
            didClose = true

            window.onbeforeunload = onbeforeunload;
            modal.append(button)
            // DEBUG: print terminal stats!
            // const quota = (println.paintedFrames / (println.missedFrames + println.paintedFrames)) * 100
            // println(`Terminal: painted=${println.paintedFrames} missed=${println.missedFrames} (${quota}%)`, true)
        }

        println("Connecting ...", true)

        // connect to the socket and send the action
        connectSocket((socket) => {
            println("Connected", true)
            socket.send(action);
            if (typeof param === 'string') {
                socket.send(param);
            }
        }, (data) => {
            println(data);
        }).then(() => {
            println("Connection closed.", true)
            close();
        }).catch(() => {
            println("Connection errored.", true)
            close();
        });
    });
})