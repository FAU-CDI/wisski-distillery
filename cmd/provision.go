package cmd

//spellchecker:words encoding json github wisski distillery internal component provision models ingredient barrel manager logging goprogram exit
import (
	"encoding/json"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/provision"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/manager"
	"github.com/FAU-CDI/wisski-distillery/pkg/logging"
	"github.com/tkw1536/goprogram/exit"
)

// Provision is the 'provision' command
var Provision wisski_distillery.Command = pv{}

type pv struct {
	PHPVersion            string `short:"p" long:"php" description:"specific php version to use for instance. See 'provision --list-php-versions' for available versions. "`
	ListPHPVersions       bool   `long:"list-php-versions" description:"List available php versions"`
	IIPServer             bool   `short:"i" long:"iip-server" description:"enable iip-server inside this instance"`
	OPCacheDevelopment    bool   `short:"o" long:"opcache-devel" description:"Include opcache development configuration"`
	Flavor                string `short:"f" long:"flavor" description:"Use specific flavor. Use '--list-flavors' to list flavors. "`
	ListFlavors           bool   `short:"l" long:"list-flavors" description:"List all known flavors"`
	ContentSecurityPolicy string `short:"c" long:"content-security-policy" description:"Setup ContentSecurityPolicy"`
	Positionals           struct {
		Slug string `positional-arg-name:"slug" description:"slug of instance to create"`
	} `positional-args:"true"`
}

var errMissingSlug = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "must provide a slug",
}

func (pv pv) AfterParse() error {
	if !pv.ListFlavors && !pv.ListPHPVersions && pv.Positionals.Slug == "" {
		return errMissingSlug
	}
	return nil
}

func (pv) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "provision",
		Description: "creates a new instance",
	}
}

var errProvisionGeneric = exit.Error{
	Message:  "unable to provision instance %s",
	ExitCode: exit.ExitGeneric,
}

// TODO: AfterParse to check instance!

func (p pv) Run(context wisski_distillery.Context) error {
	if p.ListFlavors {
		return p.listFlavors(context)
	}
	if p.ListPHPVersions {
		return p.listPHPVersions(context)
	}

	instance, err := context.Environment.Provision().Provision(context.Stderr, context.Context, provision.Flags{
		Slug:   p.Positionals.Slug,
		Flavor: p.Flavor,
		System: models.System{
			PHP:                   p.PHPVersion,
			IIPServer:             p.IIPServer,
			OpCacheDevelopment:    p.OPCacheDevelopment,
			ContentSecurityPolicy: p.ContentSecurityPolicy,
		},
	})
	if err != nil {
		return errProvisionGeneric.WithMessageF(p.Positionals.Slug).WrapError(err)
	}

	// and we're done!
	logging.LogMessage(context.Stderr, "Instance has been provisioned")
	context.Printf("URL:      %s\n", instance.URL().String())
	context.Printf("Username: %s\n", instance.DrupalUsername)
	context.Printf("Password: %s\n", instance.DrupalPassword)

	return nil
}

func (pv) listFlavors(context wisski_distillery.Context) error {
	encoder := json.NewEncoder(context.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(manager.Profiles())
}

func (pv) listPHPVersions(context wisski_distillery.Context) error {
	for _, v := range models.KnownPHPVersions() {
		if v == models.DefaultPHPVersion {
			if _, err := context.Printf("%s (default)\n", v); err != nil {
				return err
			}
		} else {
			if _, err := context.Println(v); err != nil {
				return err
			}
		}
	}
	return nil
}
