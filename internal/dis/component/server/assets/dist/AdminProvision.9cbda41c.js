var e="undefined"!=typeof globalThis?globalThis:"undefined"!=typeof self?self:"undefined"!=typeof window?window:"undefined"!=typeof global?global:{},n={},t={},o=e.parcelRequireafa4;null==o&&((o=function(e){if(e in n)return n[e].exports;if(e in t){var o=t[e];delete t[e];var r={id:e,exports:{}};return n[e]=r,o.call(r.exports,r,r.exports),r.exports}var i=new Error("Cannot find module '"+e+"'");throw i.code="MODULE_NOT_FOUND",i}).register=function(e,n){t[e]=n},e.parcelRequireafa4=o),o("8xGhL");var r=o("12vpF");async function i(e){return await new Promise(((n,t)=>{(0,r.createModal)("provision",[JSON.stringify(e)],{bufferSize:0,onClose:(o,r)=>{o?n(e.Slug):t(new Error(r??"unspecified error"))}})}))}const l=document.getElementById("provision"),a=document.getElementById("slug"),d=document.getElementById("php"),u=document.getElementById("opcacheDevelopment");l.addEventListener("submit",(e=>{e.preventDefault(),i({Slug:a.value,System:{PHP:d.value,OpCacheDevelopment:u.checked}}).then((e=>{location.href="/admin/instance/"+e})).catch((e=>{console.error(e),location.reload()}))})),l.querySelector("fieldset")?.removeAttribute("disabled");