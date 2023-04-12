package assets

import _ "embed"

// This file was automatically generated. Do not edit.

//go:embed "assets_disclaimer.txt"
var Disclaimer string

// Public holds the path to the public route 
const Public = "/⛰/"

// AssetsDefault contains assets for the 'Default' entrypoint.
var AssetsDefault = Assets{
	Scripts: `<script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.938b4407.css"><link rel="stylesheet" href="/⛰/Default.81f0a181.css">`,	
}

// AssetsUser contains assets for the 'User' entrypoint.
var AssetsUser = Assets{
	Scripts: `<script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/⛰/User.47a3b7e3.js"></script><script src="/⛰/User.924f7900.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.938b4407.css"><link rel="stylesheet" href="/⛰/User.68febbf8.css"><link rel="stylesheet" href="/⛰/User.840de3b4.css">`,	
}

// AssetsAdmin contains assets for the 'Admin' entrypoint.
var AssetsAdmin = Assets{
	Scripts: `<script nomodule="" defer src="/⛰/User.924f7900.js"></script><script type="module" src="/⛰/User.47a3b7e3.js"></script><script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/⛰/Admin.205f0180.js"></script><script src="/⛰/Admin.59fb2e50.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.938b4407.css"><link rel="stylesheet" href="/⛰/Admin.a1e05c23.css"><link rel="stylesheet" href="/⛰/User.840de3b4.css"><link rel="stylesheet" href="/⛰/User.68febbf8.css"><link rel="stylesheet" href="/⛰/Admin.6d2ae968.css">`,	
}

// AssetsAdminProvision contains assets for the 'AdminProvision' entrypoint.
var AssetsAdminProvision = Assets{
	Scripts: `<script nomodule="" defer src="/⛰/User.924f7900.js"></script><script nomodule="" defer src="/⛰/Admin.59fb2e50.js"></script><script type="module" src="/⛰/User.47a3b7e3.js"></script><script type="module" src="/⛰/Admin.205f0180.js"></script><script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/⛰/AdminProvision.3cf9e19e.js"></script><script src="/⛰/AdminProvision.d195fd59.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.938b4407.css"><link rel="stylesheet" href="/⛰/Admin.a1e05c23.css"><link rel="stylesheet" href="/⛰/User.840de3b4.css"><link rel="stylesheet" href="/⛰/User.68febbf8.css"><link rel="stylesheet" href="/⛰/Admin.6d2ae968.css"><link rel="stylesheet" href="/⛰/AdminProvision.38d394c2.css">`,	
}
