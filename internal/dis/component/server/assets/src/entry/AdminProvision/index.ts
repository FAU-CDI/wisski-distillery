import '../Admin/index.ts'
import '../Admin/index.css'

import { Provision } from '~/src/lib/remote/api'

const system = document.getElementById('system') as HTMLFormElement
const slug = document.getElementById('slug') as HTMLInputElement
const php = document.getElementById('php') as HTMLSelectElement
const phpDevelopment = document.getElementById('phpDevelopment') as HTMLInputElement
const contentSecurityPolicy = document.getElementById('contentsecuritypolicy') as HTMLInputElement
const iipserver = document.getElementById('iipserver') as HTMLInputElement
const dedicatedsql = document.getElementById('dedicatedsql') as HTMLInputElement
const dedicatedtriplestore = document.getElementById('dedicatedtriplestore') as HTMLInputElement
const ipAllowlist = document.getElementById('ipallowlist') as HTMLInputElement
const solrserver = document.getElementById('solrserver') as HTMLInputElement

// add an event handler to open the modal form!
system.addEventListener('submit', (evt) => {
  evt.preventDefault()

  const flavorElement = document.querySelector('input[name="flavor"]:checked')
  const flavor = (flavorElement instanceof HTMLInputElement) ? flavorElement.value : ''

  Provision({
    Slug: slug.value,
    Flavor: flavor,
    System: {
      PHP: php.value,
      IIPServer: iipserver.checked,
      PHPDevelopment: phpDevelopment.checked,
      ContentSecurityPolicy: contentSecurityPolicy.value,
      DedicatedSQL: dedicatedsql.checked,
      DedicatedTriplestore: dedicatedtriplestore.checked,
      IPAllowlist: ipAllowlist.value,
      SolrServer: solrserver.checked,
    },
  })
    .then(slug => {
      location.href = '/admin/instance/' + slug
    })
    .catch((e) => { console.error(e); location.reload() })
})

// enable the form!
system.querySelector('fieldset')?.removeAttribute('disabled')
