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
	Styles:  `<link rel="stylesheet" href="/⛰/Default.ad9d9c4c.css"><link rel="stylesheet" href="/⛰/Default.81f0a181.css">`,	
}

// AssetsUser contains assets for the 'User' entrypoint.
var AssetsUser = Assets{
	Scripts: `<script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/⛰/User.fce9a3e3.js"></script><script src="/⛰/User.e4c5f849.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.ad9d9c4c.css"><link rel="stylesheet" href="/⛰/User.68febbf8.css"><link rel="stylesheet" href="/⛰/User.840de3b4.css">`,	
}

// AssetsAdmin contains assets for the 'Admin' entrypoint.
var AssetsAdmin = Assets{
	Scripts: `<script nomodule="" defer src="/⛰/User.e4c5f849.js"></script><script type="module" src="/⛰/User.fce9a3e3.js"></script><script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/⛰/Admin.87f202f8.js"></script><script src="/⛰/Admin.1b10eebb.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.ad9d9c4c.css"><link rel="stylesheet" href="/⛰/Admin.a1e05c23.css"><link rel="stylesheet" href="/⛰/User.840de3b4.css"><link rel="stylesheet" href="/⛰/User.68febbf8.css"><link rel="stylesheet" href="/⛰/Admin.78d18bfa.css">`,	
}

// AssetsAdminProvision contains assets for the 'AdminProvision' entrypoint.
var AssetsAdminProvision = Assets{
	Scripts: `<script nomodule="" defer src="/⛰/User.e4c5f849.js"></script><script nomodule="" defer src="/⛰/Admin.1b10eebb.js"></script><script type="module" src="/⛰/User.fce9a3e3.js"></script><script type="module" src="/⛰/Admin.87f202f8.js"></script><script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/⛰/AdminProvision.f0586340.js"></script><script src="/⛰/AdminProvision.8e1318eb.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.ad9d9c4c.css"><link rel="stylesheet" href="/⛰/Admin.a1e05c23.css"><link rel="stylesheet" href="/⛰/User.840de3b4.css"><link rel="stylesheet" href="/⛰/User.68febbf8.css"><link rel="stylesheet" href="/⛰/Admin.78d18bfa.css"><link rel="stylesheet" href="/⛰/AdminProvision.38d394c2.css">`,	
}

// AssetsAdminRebuild contains assets for the 'AdminRebuild' entrypoint.
var AssetsAdminRebuild = Assets{
	Scripts: `<script nomodule="" defer src="/⛰/User.e4c5f849.js"></script><script nomodule="" defer src="/⛰/Admin.1b10eebb.js"></script><script type="module" src="/⛰/User.fce9a3e3.js"></script><script type="module" src="/⛰/Admin.87f202f8.js"></script><script type="module" src="/⛰/Default.38d394c2.js"></script><script src="/⛰/Default.38d394c2.js" nomodule="" defer></script><script type="module" src="/⛰/AdminRebuild.7523d442.js"></script><script src="/⛰/AdminRebuild.d7e40392.js" nomodule="" defer></script>`,
	Styles:  `<link rel="stylesheet" href="/⛰/Default.ad9d9c4c.css"><link rel="stylesheet" href="/⛰/Admin.a1e05c23.css"><link rel="stylesheet" href="/⛰/User.840de3b4.css"><link rel="stylesheet" href="/⛰/User.68febbf8.css"><link rel="stylesheet" href="/⛰/Admin.78d18bfa.css"><link rel="stylesheet" href="/⛰/AdminRebuild.38d394c2.css">`,	
}
