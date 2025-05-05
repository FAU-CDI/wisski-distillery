package cmd

//spellchecker:words github wisski distillery internal component server assets
import (
	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
)

// License is the 'wdcli license' command.
//
// The license command prints to standard output legal notices about the wdcli program.
var License wisski_distillery.Command = license{}

type license struct{}

func (license) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: false,
		},
		Command:     "license",
		Description: "print licensing information about wdcli and exit",
	}
}

func (license) AfterParse() error {
	return nil
}

func (license) Run(context wisski_distillery.Context) error {
	_, _ = context.Printf(stringLicenseInfo, wisski_distillery.License, cli.LegalNotices, assets.Disclaimer)
	return nil
}

const stringLicenseInfo = `
wdcli -- WissKI Distillery Command Line Utility
https://github.com/FAU-CDI/wisski-distillery

================================================================================
wdcli is licensed under the terms of the AGPL Version 3.0 License:

%s
================================================================================

Furthermore, this executable may include code from the following projects:
%s

================================================================================

Finally, the web frontend may contain additional code.
%s
`
