import './index.css'

import { discard } from '~/src/lib/discard'

document.querySelectorAll('.copy').forEach((elem: Element) => {
  elem.addEventListener('click', () => {
    // Check if the clipboard api is supported
    // eslint-disable-next-line @typescript-eslint/strict-boolean-expressions
    if (!navigator.clipboard) return

    discard(navigator.clipboard.writeText((elem as HTMLElement).innerText))
  })
})
