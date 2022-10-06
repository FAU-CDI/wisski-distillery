import '../global.js';
import './index.css';

const types = {
    "date": (element) => {
        return (new Date(element.innerText)).toISOString()
    },
    "path": (element) => {
        const text = element.innerText.split("/");
        return text[text.length - 1];
    },
    "pathbuilders": (element) => {
        const pathbuilders = window.pathbuilders; // read from context!
        const wrapper = document.createElement("span");

        let found_one = false
        Object.keys(pathbuilders).forEach(name => {
            found_one = true

            const filename = name + ".xml"
            const data = pathbuilders[name]
            const mime = "application/xml"
            wrapper.append(make_download_link(filename, name, data, mime))
            wrapper.append(document.createTextNode(" "))
        })

        if (!found_one) return '(none)';

        const small = document.createElement('small')
        small.append(document.createTextNode("(click to download)"))
        wrapper.append(small)
        
        return wrapper
    }
}

const make_download_link = (filename, title, content, type) => {
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
    link.append(document.createTextNode(title))

    return link
}

Object.keys(types).forEach(key => {
    const f = types[key];
    const elements = document.querySelectorAll("code." + key)
    elements.forEach(element => {
        const newElement = f(element)
        if (typeof newElement === 'string') {
            element.innerHTML = ""
            element.appendChild(document.createTextNode(newElement))
            return
        }

        element.parentNode.replaceChild(newElement, element)
    })
})