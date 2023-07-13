import { Parcel } from '@parcel/core'
import { mkdir, rm, writeFile, readFile, unlink } from 'fs/promises'
import { join } from 'path'
import { parse as parseHTML } from 'node-html-parser'
import { spawnSync } from 'child_process'

//
// PARAMETERS
//

const ENTRYPOINTS = process.argv.slice(2)
const ENTRY_DIR = join('.', '.entry-cache') // directory to place entries into
const DIST_DIR = join('.', 'dist')
const PUBLIC_DIR = '/⛰/' // mountain's don't move, and neither do static files

const DEST_PACKAGE = process.env.GOPACKAGE ?? 'static'
const DEST_DISCLAIMER = (() => {
  const source = (process.env.GOFILE ?? 'assets.go')
  const base = source.substring(0, source.length - '.go'.length)
  return base + '_disclaimer.txt'
})()
const DEST_FILE = (() => {
  const source = (process.env.GOFILE ?? 'assets.go')
  const base = source.substring(0, source.length - '.go'.length)
  return base + '_dist.go'
})()

//
// PREPARE DIRECTORIES
//

process.stdout.write('Preparing directories ...')
await Promise.all([
  mkdir(ENTRY_DIR, { recursive: true }),
  rm(DIST_DIR, { recursive: true, force: true })
])
console.log(' Done.')

//
// Write the disclaimer
//

process.stdout.write('Generating legal disclaimer ...')

const disclaimer = await new Promise((resolve, reject) => {
  const child = spawnSync('yarn', ['licenses', 'generate-disclaimer'], { encoding: 'utf8' })
  if (child.error) {
    reject(child.stderr)
    return
  }

  resolve(child.stdout)
})

console.log(' Done.')

process.stdout.write(`Writing ${DEST_DISCLAIMER} ...`)
await writeFile(DEST_DISCLAIMER, disclaimer)
console.log(' Done.')

//
// WRITE ENTRY POINTS
//

process.stdout.write('Collecting entry points ')
const entries = await Promise.all(ENTRYPOINTS.map(async (name) => {
  const entry = {
    name,
    bundleName: name + '.html',
    src: join(ENTRY_DIR, name + '.html')
  }

  const content = `
<script type='module' src='../src/base/index.ts'></script>
<script type='module' src='../src/entry/${name}/index.ts'></script>
<link rel='stylesheet' href='../src/entry/${name}/index.css'>
`
  await writeFile(entry.src, content)

  process.stdout.write('.')
  return entry
}))
console.log(' Done.')

//
// BUNDLEING
//

process.stdout.write('Bundleing assets ...')
const bundler = new Parcel({
  entries: entries.map(e => e.src),
  defaultConfig: '@parcel/config-default',
  shouldDisableCache: true,
  shouldContentHash: true,
  defaultTargetOptions: {
    shouldOptimize: true,
    shouldScopeHoist: true,
    sourceMaps: false,
    distDir: DIST_DIR,
    publicUrl: PUBLIC_DIR,
    engines: {
      browsers: 'defaults'
    }
  }
})
const { bundleGraph } = await bundler.run()
console.log(' Done.')

//
// FIND ASSETS IN OUTPUT
//

process.stdout.write('Find Assets in Output ')
const bundles = bundleGraph.getBundles()
const assets = await Promise.all(entries.map(async (entry) => {
  const mainBundle = bundles.find(b => b.name === entry.bundleName)
  if (mainBundle === undefined) throw new Error('Unable to find bundle for ' + entry.name)

  // read, then delete the generated output file
  const { filePath } = mainBundle
  const html = parseHTML(await readFile(filePath))
  await unlink(filePath)

  const scripts = html.querySelectorAll('script').map(script => script.outerHTML).join('')
  const links = html.querySelectorAll('link').map(link => link.outerHTML).join('')

  process.stdout.write('.')
  return { ...entry, scripts, links }
}))
console.log(' Done.')

//
// GENERATE GO
//

process.stdout.write(`Writing ${DEST_FILE} ...`)
const goAssets = assets.map(({ name, scripts, links }) => {
  return `
// Assets${name} contains assets for the '${name}' entrypoint.
var Assets${name} = Assets{
\tScripts: \`${scripts}\`,
\tStyles:  \`${links}\`,\t
}`.trim()
}).join('\n\n')
const goSource = `package ${DEST_PACKAGE}

import _ "embed"

// This file was automatically generated. Do not edit.

//go:embed ${JSON.stringify(DEST_DISCLAIMER)}
var Disclaimer string

// Public holds the path to the public route 
const Public = ${JSON.stringify(PUBLIC_DIR)}

${goAssets}
`

await writeFile(DEST_FILE, goSource)
console.log(' Done.')
