// setup highlighting
import '~/src/lib/highlight'

// setup remote actions
import setup from '~/src/lib/remote'

// csv table export
import '~/src/lib/csv-export'

// include the user styles!
import '../User/index.ts'
import '../User/index.css'

// highlight everything
import 'highlight.js/styles/default.css'
import highlightJs from 'highlight.js'
setup()
highlightJs.highlightAll()
