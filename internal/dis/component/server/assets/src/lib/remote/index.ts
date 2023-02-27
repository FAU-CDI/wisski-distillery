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

export default function setup() {
    const remote_action = document.getElementsByClassName('remote-action')
    Array.from(remote_action).forEach((element) => {
        const action = element.getAttribute('data-action') as string;
        const reload = element.getAttribute('data-force-reload');
        const param = element.getAttribute('data-param') as string | undefined;
        
        const confirmElementName = element.getAttribute('data-confirm-param');
        const confirmElement = (confirmElementName ? document.querySelector(confirmElementName) : null) as HTMLInputElement | null;
    
        const bufferSize = (function () {
            const number = parseInt(element.getAttribute('data-buffer') ?? "", 10) ?? 0;
            return (isFinite(number) && number > 0) ? number : 0;
        })()
    
        const validate = function() {
            if (!confirmElement) return true
            return confirmElement.value === param;
        }
    
        if (confirmElement) {
            const runValidation = () => {
                if (validate()) {
                    element.removeAttribute('disabled')
                } else {
                    element.setAttribute('disabled', 'disabled')
                }
            }
            confirmElement.addEventListener('change', runValidation)
            runValidation()
        }

        let onClose: (success: boolean) => void | null;
        if (typeof reload === 'string') {
            onClose = () => {
                if (reload === '') location.reload();
                else location.href = reload;
            }
        }
    
        element.addEventListener('click', function (ev) {
            ev.preventDefault();
    
            // do nothing if the validation fails
            if (!validate()) return;

            // create a modal dialog
            const params = (typeof param === 'string') ? [param] : [];
            createModal(action, params, {
                onClose: onClose, 
                bufferSize: bufferSize,
            });
        });
    })
}


type ModalOptions = {
    bufferSize: number;
    onClose: (success: boolean) => void 
}
export function createModal(action: string, params: string[], opts: Partial<ModalOptions>) {
    // create a modal dialog and append it to the body
    const modal = document.createElement("div")
    modal.className = "modal-terminal"
    document.body.append(modal)

    // create a <pre> to write stuff into
    const target = document.createElement("pre")
    const println = makeTextBuffer(target, modal, opts.bufferSize ?? 1000)
    modal.append(target)
    
    // create a button to eventually close everything
    const button = document.createElement("button")
    button.className = "pure-button pure-button-success"
    button.append(typeof opts?.onClose === 'function' ? "Close & Finish" : "Close")
    let success = false;
    button.addEventListener('click', function (event) {
        event.preventDefault();

        if (typeof opts?.onClose === 'function') {
            button.setAttribute('disabled', 'disabled')
            target.innerHTML = 'Finishing up ...'
            opts.onClose(success)
            return;
        }

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
        socket.send(action)
        params.forEach(p => socket.send(p))
    }, (data) => {
        println(data);
    }).then(() => {
        success = true
        println("Connection closed.", true)
        close();
    }).catch(() => {
        success = false
        println("Connection errored.", true)
        close();
    });
}