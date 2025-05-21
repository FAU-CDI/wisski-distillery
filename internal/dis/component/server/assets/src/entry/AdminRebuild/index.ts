import '../Admin/index.ts'
import '../Admin/index.css'

import { Rebuild } from '~/src/lib/remote/api'

const system = document.getElementById('system') as HTMLFormElement
const slug = document.getElementById('slug') as HTMLInputElement
const php = document.getElementById('php') as HTMLSelectElement
const phpDevelopment = document.getElementById('phpDevelopment') as HTMLInputElement
const contentSecurityPolicy = document.getElementById('contentsecuritypolicy') as HTMLInputElement
const iipserver = document.getElementById('iipserver') as HTMLInputElement

// add an event handler to open the modal form!
system.addEventListener('submit', (evt) => {
  evt.preventDefault()

  Rebuild(slug.value, { PHP: php.value, IIPServer: iipserver.checked, PHPDevelopment: phpDevelopment.checked, ContentSecurityPolicy: contentSecurityPolicy.value })
    .then(slug => {
      location.href = '/admin/instance/' + slug
    })
    .catch((e) => { console.error(e); location.reload() })
})

// enable the form!
system.querySelector('fieldset')?.removeAttribute('disabled')
