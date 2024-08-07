import { Result } from "../socketapi/pow_client"
import './index.css'
import LocalSession from './local'

type Print = ((text: string, flush?: boolean) => void) & {
  paintedFrames: number
  missedFrames: number
}

const NEW_LINE = '\n'
const NEW_LINE_LENGTH = NEW_LINE.length

/**
 * trimLines trims buffer so that it contains as most count lines
 */
function trimLines (buffer: string, lines: number): string {
  if (lines <= 0 || isNaN(lines) || !isFinite(lines)) return buffer

  let count = 0
  let index = buffer.length

  // while we still have sufficient space
  while (count < lines) {
    // get the next start of the line
    index = buffer.lastIndexOf(NEW_LINE, index - 1)
    if (index === -1) {
      return buffer
    }

    // increase the count
    count++
  }

  return buffer.substring(index + NEW_LINE_LENGTH)
}

/**
 * makeTextBuffer returns a println() function that efficiently writes text into target, and keeps at most size elements in the traceback.
 * scrollContainer is used to scroll on every painted update.
 */
function makeTextBuffer (target: HTMLElement, scrollContainer: HTMLElement, size: number): Print {
  let lastAnimationFrame: number | null = null // last scheduled animation frame

  // text buffer
  let buffer: string = ''
  const paint = (): void => {
    print.paintedFrames++
    target.innerText = buffer
    scrollContainer.scrollTop = scrollContainer.scrollHeight
    lastAnimationFrame = null
  }

  const print = (text: string, flush?: boolean): void => {
    // add text to the buffer and normalize
    buffer += text.replace(/^\s*[\r\n]/gm, '\r\n')

    // trim the buffer to the specified number of lines
    buffer = trimLines(buffer, size)
    // and update the browser in the next animation frame
    if (lastAnimationFrame !== null) {
      print.missedFrames++
      window.cancelAnimationFrame(lastAnimationFrame)
    }

    // force a repaint!
    if (flush === true) return paint()

    // schedule an animation frame
    lastAnimationFrame = window.requestAnimationFrame(paint)
  }
  print.paintedFrames = 0
  print.missedFrames = 0

  return print
}

export default function setup (): void {
  const remoteAction = document.getElementsByClassName('remote-action')
  Array.from(remoteAction).forEach((element) => {
    const action = element.getAttribute('data-action') as string
    const reload = element.getAttribute('data-force-reload')
    const param = element.getAttribute('data-param') as string | undefined

    const confirmElementName = element.getAttribute('data-confirm-param')
    const confirmElement = typeof confirmElementName === 'string' ? document.querySelector(confirmElementName) : null

    const getConfirmValue = (): string | null => {
      if (confirmElement === null) {
        console.warn('data-confirm-param was not found')
        return null
      }
      if (!('value' in confirmElement)) {
        return null
      }
      const value = confirmElement.value
      if (value === null || (typeof value !== 'string')) {
        return null
      }

      return value
    }

    const bufferSize = (function () {
      const number = parseInt(element.getAttribute('data-buffer') ?? '', 10) ?? 0
      return (isFinite(number) && number > 0) ? number : 0
    })()

    const validate = function (): boolean {
      const confirmValue = getConfirmValue()
      if (confirmValue === null) return true
      return confirmValue === param
    }

    if (confirmElement !== null) {
      const runValidation = (): void => {
        if (validate()) {
          element.removeAttribute('disabled')
        } else {
          element.setAttribute('disabled', 'disabled')
        }
      }
      confirmElement.addEventListener('change', runValidation)
      runValidation()
    }

    let onClose: ((success: boolean, data: any) => void) | undefined
    if (typeof reload === 'string') {
      onClose = () => {
        location.href = reload === '' ? location.href : reload
      }
    }

    element.addEventListener('click', function (ev) {
      ev.preventDefault()

      // do nothing if the validation fails
      if (!validate()) return

      // create a modal dialog
      const params = (typeof param === 'string') ? [param] : []
      createModal(action, params, {
        onClose,
        bufferSize
      })
    })
  })
}

interface ModalOptions {
  bufferSize: number
  onClose: ((success: true, data: any) => void) & ((success: false, message: string) => void)
}
export function createModal (action: string, params: string[], opts: Partial<ModalOptions>): void {
  // create a modal dialog and append it to the body
  const modal = document.createElement('div')
  modal.className = 'modal-terminal'
  document.body.append(modal)

  // create a <pre> to write stuff into
  const target = document.createElement('pre')
  const print = makeTextBuffer(target, modal, opts.bufferSize ?? 1000)
  modal.append(target)

  // create a button to eventually close everything
  const finishButton = document.createElement('button')
  finishButton.className = 'pure-button pure-button-success'
  finishButton.append(typeof opts?.onClose === 'function' ? 'Close & Finish' : 'Close')

  let result: Result = { status: 'rejected', reason: 'Nothing happened' }
  finishButton.addEventListener('click', (event) => {
    event.preventDefault()

    if (typeof opts?.onClose === 'function') {
      finishButton.setAttribute('disabled', 'disabled')
      target.innerHTML = 'Finishing up ...'
      if (result.status === 'fulfilled') {
        opts.onClose(true, result.value)
      } else {
        opts.onClose(false, result.reason ?? 'unknown error')
      }
      return
    }

    modal.parentNode?.removeChild(modal)
  })

  const cancelButton = document.createElement('button')
  cancelButton.className = 'pure-button pure-button-danger'
  cancelButton.setAttribute('disabled', 'disabled')
  cancelButton.append('Cancel')
  modal.append(cancelButton)

  const onbeforeunload = window.onbeforeunload
  window.onbeforeunload = () => 'A remote session is in progress. Are you sure you want to leave?'

  // when closing, add a button to the modal!
  const close = (message: Result): void => {
    result = message

    if (result.status === 'fulfilled') {
      print('Process completed successfully.\n', true)
    } else {
      print('Process reported error: ' + (result.reason ?? 'unknown error') + '\n', true)
    }

    window.onbeforeunload = onbeforeunload

    modal.removeChild(cancelButton)
    modal.append(finishButton)

    const quota = (print.paintedFrames / (print.missedFrames + print.paintedFrames)) * 100
    console.debug(`Result:`, result)
    console.debug(`Terminal: painted=${print.paintedFrames} missed=${print.missedFrames} (${quota}%)`)
  }

  print('Connecting ...', true)

  // connect to the socket and send the action
  const session = new LocalSession({
    call: action,
    params
  })

  session.beforeCall = function () {
    cancelButton.removeAttribute('disabled')
    cancelButton.addEventListener('click', (event) => {
      event.preventDefault()

      print('^C\n', true)
      this.cancel()
    })
    print(' Connected.\n', true)
  }
  session.onLogLine = print

  session.connect()
    .then(() => session.closeInput()) // for now none of our sessions actually have input
    .then(() => session.wait())
    .then((result) => {
      close(result.result)
    })
    .catch((err) => {
      console.error(err)
      close({ status: 'rejected', reason: 'connection closed unexpectedly' })
    })
}
