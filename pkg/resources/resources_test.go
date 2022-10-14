package resources

import (
	"fmt"
	"strings"
)

func ExampleParse() {
	resources := Parse(strings.NewReader(`
	<html>
		<head>
			<link rel="stylesheet" href="/some/sheet1.css">
			<link rel="stylesheet" href="/some/sheet2.css">
		</head>
		<body>
			<script type="module" src="/some/module1.js"></script>
			<script type="module" src="/some/module2.js"></script>
			<script src="/some/nonmodule1.js"></script>
			<script src="/some/nonmodule2.js"></script>
		</body>
	</html>
	`))

	var builder strings.Builder
	builder.WriteString("css: ")
	resources.WriteCSS(&builder)

	builder.WriteString("\njs: ")
	resources.WriteJS(&builder)
	fmt.Println(builder.String())

	// Output: css: <link rel=stylesheet href="/some/sheet1.css"><link rel=stylesheet href="/some/sheet2.css">
	// js: <script type=module src="/some/module1.js"><script type=module src="/some/module2.js"><script src="/some/nonmodule1.js"><script src="/some/nonmodule2.js">
}
