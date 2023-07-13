import './index.css'

/** Adapted from http://blog.parkermoore.de/2014/08/01/header-anchor-links-in-vanilla-javascript-for-github-pages-and-jekyll/ */
const anchorForId = (id: string): HTMLAnchorElement => {
  const anchor = document.createElement('a')
  anchor.className = 'header-link'
  anchor.href = '#' + id
  anchor.innerHTML = '#'
  return anchor
}

const linkifyAnchors = (level: number): void => {
  const headers = document.getElementsByTagName('h' + level.toString())
  Array.from(headers).forEach((header) => {
    if (typeof header.id === 'undefined' || header.id === '') return
    header.appendChild(anchorForId(header.id))
  })
}

// linkify all the anchors from 1 ... 6
(new Array(6)).fill(0).forEach((_, i) => linkifyAnchors(i + 1))
