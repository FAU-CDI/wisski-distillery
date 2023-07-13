document.querySelectorAll('span').forEach((elem: Element) => {
  if (!elem.hasAttribute('data-reveal')) return

  addReveal(elem as HTMLSpanElement, 10000)
})

export function addReveal (span: HTMLSpanElement, hideDelay: number): void {
  const content = span.getAttribute('data-reveal') ?? '(no content)'

  let isHidden = true

  // handler to hide the element
  const hide = (): void => {
    isHidden = true
    span.innerText = '(click to reveal)'
  }
  hide()

  const reveal = (): void => {
    isHidden = false
    const code = document.createElement('code')
    code.append(content)
    code.addEventListener('click', (evt) => {
      evt.preventDefault()

      // Check if the clipboard api is supported
      // eslint-disable-next-line @typescript-eslint/strict-boolean-expressions
      if (!navigator.clipboard) return

      navigator.clipboard.writeText(content).then(() => {}).catch(console.error.bind(console))
    })
    code.style.userSelect = 'all'

    span.innerHTML = ''
    span.append(code)
  }

  span.addEventListener('click', (evt) => {
    evt.preventDefault()

    if (!isHidden) return
    reveal()
    setTimeout(hide, hideDelay) // hide again after 1 second
  })
}
