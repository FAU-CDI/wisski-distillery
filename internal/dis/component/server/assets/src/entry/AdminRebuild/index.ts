import '../Admin/index.ts'
import '../Admin/index.css'

import { Rebuild } from '~/src/lib/remote/api'

const rebuild = document.getElementById('rebuild') as HTMLFormElement
const slug = document.getElementById('slug') as HTMLInputElement
const php = document.getElementById('php') as HTMLSelectElement
const opcacheDevelopment = document.getElementById('opcacheDevelopment') as HTMLInputElement
const contentSecurityPolicy = document.getElementById('contentsecuritypolicy') as HTMLInputElement

// add an event handler to open the modal form!
rebuild.addEventListener('submit', (evt) => {
  evt.preventDefault()

  Rebuild(slug.value, { PHP: php.value, OpCacheDevelopment: opcacheDevelopment.checked, ContentSecurityPolicy: contentSecurityPolicy.value })
    .then(slug => {
      location.href = '/admin/instance/' + slug
    })
    .catch((e) => { console.error(e); location.reload() })
})

// enable the form!
rebuild.querySelector('fieldset')?.removeAttribute('disabled')
