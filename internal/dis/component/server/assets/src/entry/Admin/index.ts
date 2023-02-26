import "~/src/lib/remote"
import "~/src/lib/highlight"

// include the user styles!
import "../User/index.ts"
import "../User/index.css"

// highlight everything
import "highlight.js/styles/default.css"
import highlightJs from "highlight.js"
highlightJs.highlightAll();