import "../Admin/index.ts"
import "../Admin/index.css"

import { Provision } from "~/src/lib/remote/api"

const provision = document.getElementById("provision") as HTMLFormElement;
const slug = document.getElementById("slug") as HTMLInputElement;
const php = document.getElementById("php") as HTMLSelectElement;
const opcacheDevelopment = document.getElementById("opcacheDevelopment") as HTMLInputElement;

// add an event handler to open the modal form!
provision.addEventListener('submit', (evt) => {
    evt.preventDefault();

    Provision({ Slug: slug.value, System: { PHP: php.value, OpCacheDevelopment: opcacheDevelopment.checked } })
        .then(slug => {
            location.href = "/admin/instance/" + slug;
        })
        .catch((e) => {console.error(e); location.reload()});
})

// enable the form!
provision.querySelector('fieldset')?.removeAttribute('disabled');

