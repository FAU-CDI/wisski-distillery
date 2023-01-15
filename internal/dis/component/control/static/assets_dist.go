package static

import _ "embed"

// This file was automatically generated. Do not edit.

//go:embed "assets_disclaimer.txt"
var AssetsDisclaimer string

// AssetsDefault contains assets for the 'Default' entrypoint.
var AssetsDefault = Assets{
	Scripts: `<script type="module" src="/static/Default.38d394c2.js"></script><script src="/static/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/static/Default.38d394c2.js"></script><script src="/static/Default.38d394c2.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/static/Default.db26a303.css"><link rel="stylesheet" href="/static/Default.f9675eae.css">`,	
}

// AssetsUser contains assets for the 'User' entrypoint.
var AssetsUser = Assets{
	Scripts: `<script type="module" src="/static/Default.38d394c2.js"></script><script src="/static/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/static/User.e0367d79.js"></script><script src="/static/User.b2f9a57c.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/static/Default.db26a303.css"><link rel="stylesheet" href="/static/User.68febbf8.css"><link rel="stylesheet" href="/static/User.840de3b4.css">`,	
}

// AssetsAdmin contains assets for the 'Admin' entrypoint.
var AssetsAdmin = Assets{
	Scripts: `<script nomodule="" defer src="/static/User.b2f9a57c.js"></script><script type="module" src="/static/User.e0367d79.js"></script><script type="module" src="/static/Default.38d394c2.js"></script><script src="/static/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/static/Admin.4ca3cb6f.js"></script><script src="/static/Admin.9750ba9c.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/static/Default.db26a303.css"><link rel="stylesheet" href="/static/Admin.6d59e220.css"><link rel="stylesheet" href="/static/User.840de3b4.css"><link rel="stylesheet" href="/static/User.68febbf8.css"><link rel="stylesheet" href="/static/Admin.6d2ae968.css">`,	
}
