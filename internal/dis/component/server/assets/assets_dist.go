package assets

import _ "embed"

// This file was automatically generated. Do not edit.

//go:embed "assets_disclaimer.txt"
var Disclaimer string

// Public holds the path to the public route 
const Public = "/this-is-fine/"

// AssetsDefault contains assets for the 'Default' entrypoint.
var AssetsDefault = Assets{
	Scripts: `<script type="module" src="/this-is-fine/Default.38d394c2.js"></script><script src="/this-is-fine/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/this-is-fine/Default.38d394c2.js"></script><script src="/this-is-fine/Default.38d394c2.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/this-is-fine/Default.938b4407.css"><link rel="stylesheet" href="/this-is-fine/Default.81f0a181.css">`,	
}

// AssetsUser contains assets for the 'User' entrypoint.
var AssetsUser = Assets{
	Scripts: `<script type="module" src="/this-is-fine/Default.38d394c2.js"></script><script src="/this-is-fine/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/this-is-fine/User.e0367d79.js"></script><script src="/this-is-fine/User.b2f9a57c.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/this-is-fine/Default.938b4407.css"><link rel="stylesheet" href="/this-is-fine/User.68febbf8.css"><link rel="stylesheet" href="/this-is-fine/User.840de3b4.css">`,	
}

// AssetsAdmin contains assets for the 'Admin' entrypoint.
var AssetsAdmin = Assets{
	Scripts: `<script nomodule="" defer src="/this-is-fine/User.b2f9a57c.js"></script><script type="module" src="/this-is-fine/User.e0367d79.js"></script><script type="module" src="/this-is-fine/Default.38d394c2.js"></script><script src="/this-is-fine/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/this-is-fine/Admin.1f7e8ff8.js"></script><script src="/this-is-fine/Admin.9f3b0d76.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/this-is-fine/Default.938b4407.css"><link rel="stylesheet" href="/this-is-fine/Admin.a1e05c23.css"><link rel="stylesheet" href="/this-is-fine/User.840de3b4.css"><link rel="stylesheet" href="/this-is-fine/User.68febbf8.css"><link rel="stylesheet" href="/this-is-fine/Admin.6d2ae968.css">`,	
}

// AssetsAdminProvision contains assets for the 'AdminProvision' entrypoint.
var AssetsAdminProvision = Assets{
	Scripts: `<script nomodule="" defer src="/this-is-fine/User.b2f9a57c.js"></script><script nomodule="" defer src="/this-is-fine/Admin.9f3b0d76.js"></script><script type="module" src="/this-is-fine/User.e0367d79.js"></script><script type="module" src="/this-is-fine/Admin.1f7e8ff8.js"></script><script type="module" src="/this-is-fine/Default.38d394c2.js"></script><script src="/this-is-fine/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/this-is-fine/AdminProvision.3cf9e19e.js"></script><script src="/this-is-fine/AdminProvision.d195fd59.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/this-is-fine/Default.938b4407.css"><link rel="stylesheet" href="/this-is-fine/Admin.a1e05c23.css"><link rel="stylesheet" href="/this-is-fine/User.840de3b4.css"><link rel="stylesheet" href="/this-is-fine/User.68febbf8.css"><link rel="stylesheet" href="/this-is-fine/Admin.6d2ae968.css"><link rel="stylesheet" href="/this-is-fine/AdminProvision.38d394c2.css">`,	
}