import "./index.css"

document.querySelectorAll('.copy').forEach((elem: Element) => {
    elem.addEventListener('click', () => {
        if (!navigator.clipboard) return;
         navigator.clipboard.writeText((elem as HTMLElement).innerText);
    })
})