document.querySelectorAll('span').forEach((elem: Element) => {
    if (!elem.hasAttribute('data-reveal')) return

    addReveal(elem as HTMLSpanElement, 10000);
})

export function addReveal(span: HTMLSpanElement, hideDelay: number) {
    const content = span.getAttribute("data-reveal") ?? '(no content)'

    let isHidden = true

    // handler to hide the element
    const hide = () => {
        isHidden = true
        span.innerText = "(click to reveal)"
    }
    hide()
    
    const reveal = () => {
        isHidden = false
        const code = document.createElement('code')
        code.append(content)
        code.addEventListener('click', (evt) => {
            evt.preventDefault()
        
            if (!navigator.clipboard) return
            navigator.clipboard.writeText(content)
        })
        code.style.userSelect = "all";

        span.innerHTML = ""
        span.append(code)
    }

    span.addEventListener("click", (evt) => {
        evt.preventDefault()


        if (!isHidden) return
        reveal()
        setTimeout(hide, hideDelay) // hide again after 1 second
    })
} 