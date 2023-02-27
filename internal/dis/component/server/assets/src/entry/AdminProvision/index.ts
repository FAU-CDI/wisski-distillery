import "../Admin/index.ts"
import "../Admin/index.css"

import { createModal } from "~/src/lib/remote"

const provision = document.getElementById("provision") as HTMLFormElement;
const slug = document.getElementById("slug") as HTMLInputElement;

// add an event handler to open the modal form!
provision.addEventListener('submit', (evt) => {
    evt.preventDefault();

    // flags used to create the server
    const flags = { Slug: slug.value };

    // open a modal to provision a new instance
    createModal("provision", [JSON.stringify(flags)], {
        bufferSize: 0,
        onClose: (success: boolean) => {
            if (success) {
                location.href = "/admin/instance/" + flags.Slug
            } else {
                location.reload();
            }
        },
    })

})

// enable the form!
provision.querySelector('fieldset')?.removeAttribute('disabled');

