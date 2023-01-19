import dayjs from "dayjs"
const types: Record<string, (element: HTMLElement) => HTMLElement | string> = {
    "date": (element) => {
        const value =  dayjs(element.innerText);
        const text = value.format('YYYY-MM-DD HH:mm:ss ([UTC]Z)')

        // if the date is the zero date, then it is assumed to be invalid
        if (value.unix() === 0) {
            const code = document.createElement('code')
            code.style.color = 'gray'
            code.append(text)
            return code
        }
        return text
    },
    "path": (element) => {
        const text = element.innerText.split("/");
        return text[text.length - 1];
    },
    "pathbuilder": (element) => {
        // create a link and get the blob
        const filename = (element.getAttribute('data-name') ?? 'pathbuilder') + ".xml"
        const [link, blob] = make_download_link(filename, element.innerText, "application/xml")

        link.className = "pure-button"
        const title = filename + ' (' + blob.size + ' Bytes)';
        link.append(title)
        return link
    }
}

const make_download_link = (filename: string, content: string, type: string): [HTMLAnchorElement, Blob] => {
    const blob = new Blob(
        [content],
        {
            type: type ?? "text/plain"
        }
    );

    const link = document.createElement("a")
    link.target = "_blank"
    link.download = filename
    link.href = URL.createObjectURL(blob)

    return [link, blob]
}

Object.keys(types).forEach(key => {
    const f = types[key];
    const elements = document.querySelectorAll("code." + key) as NodeListOf<HTMLElement>
    elements.forEach(element => {
        const newElement = f(element)
        if (typeof newElement === 'string') {
            element.innerHTML = ""
            element.appendChild(document.createTextNode(newElement))
            return
        }

        element.parentNode!.replaceChild(newElement, element)
    })
})