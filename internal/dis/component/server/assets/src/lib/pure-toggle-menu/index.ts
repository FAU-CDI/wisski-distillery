import './index.css'

const WINDOW_CHANGE_EVENT = ('onorientationchange' in window) ? 'orientationchange' : 'resize'

document.querySelectorAll('.pure-toggle-menu').forEach((menu) => {
  const toggle = menu.querySelector('.toggle')
  if (toggle == null) {
    console.warn("'.pure-toggle-menu' without '.toggle'")
    return
  }

  const toggleMenu = (): void => {
    menu.classList.toggle('closed')
    toggle.classList.toggle('x')
  }

  toggle.addEventListener('click', (e) => {
    e.preventDefault()
    toggleMenu()
  })

  window.addEventListener(WINDOW_CHANGE_EVENT, () => {
    if (menu.classList.contains('closed')) return
    toggleMenu()
  })
})
