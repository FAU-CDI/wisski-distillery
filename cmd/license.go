package cmd

//spellchecker:words github wisski distillery internal component server assets
import (
	"fmt"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/notices"
	"github.com/spf13/cobra"
)

func NewLicenseCommand() *cobra.Command {
	impl := new(license)

	cmd := &cobra.Command{
		Use:   "license",
		Short: "print licensing information about wdcli and exit",
		Args:  cobra.NoArgs,
		RunE:  impl.Exec,
	}

	return cmd
}

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

func (license) Exec(cmd *cobra.Command, args []string) error {
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), stringLicenseInfo, wisski_distillery.License, notices.LegalNotices, assets.Disclaimer)
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
