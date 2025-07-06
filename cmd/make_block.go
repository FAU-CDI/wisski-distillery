package cmd

//spellchecker:words html template time github wisski distillery internal wdlog ingredient extras goprogram exit
import (
	"fmt"
	"html/template"
	"io"
	"time"

	wisski_distillery "github.com/FAU-CDI/wisski-distillery"
	"github.com/FAU-CDI/wisski-distillery/internal/cli"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
	"go.tkw01536.de/goprogram/exit"
)

// MakeBlock is the 'make_block' command.
var MakeBlock wisski_distillery.Command = makeBlock{}

type makeBlock struct {
	Title  string `description:"title of block to create"           long:"title"  short:"t"`
	Region string `description:"optional region to assign block to" long:"region" short:"r"`
	Footer bool   `description:"create block in the footer region"  long:"footer" short:"f"`

	Positionals struct {
		Slug string `description:"slug of instance to create legal block for" positional-arg-name:"SLUG" required:"1-1"`
	} `positional-args:"true"`
}

var errFooterAndRegion = exit.NewErrorWithCode("`--footer` and `--region` provided", exit.ExitCommandArguments)

func (mb makeBlock) AfterParse() error {
	if mb.Region != "" && mb.Footer {
		return errFooterAndRegion
	}
	return nil
}

func (makeBlock) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "make_block",
		Description: "Creates a block with html content provided on stdin",
	}
}

var (
	errBlocksGeneric      = exit.NewErrorWithCode("unable to create block", exit.ExitGeneric)
	errBlocksFooterFailed = exit.NewErrorWithCode("unable to determine footer block", exit.ExitGeneric)
	errBlocksNoFooter     = exit.NewErrorWithCode("no footer known for region", exit.ExitGeneric)
	errBlocksNoContent    = exit.NewErrorWithCode("unable to read content from standard input", exit.ExitCommandArguments)
)

func (mb makeBlock) Run(context wisski_distillery.Context) error {
	// get the wisski
	instance, err := context.Environment.Instances().WissKI(context.Context, mb.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errPathbuildersWissKI, err)
	}

	// get the footer (if any)
	if mb.Footer {
		wdlog.Of(context.Context).Info("checking for footer")
		region, err := instance.Blocks().GetFooterRegion(context.Context, nil)
		if err != nil {
			return fmt.Errorf("%w: %w", errBlocksFooterFailed, err)
		}
		if region == "" {
			return errBlocksNoFooter
		}
	}

	id := ""
	if mb.Region != "" {
		id = fmt.Sprintf("block-auto-%d", time.Now().Unix())
	}

	// read the content
	content, err := io.ReadAll(context.Stdin)
	if err != nil {
		return fmt.Errorf("%w: %w", errBlocksNoContent, err)
	}

	{
		err := instance.Blocks().Create(context.Context, nil, extras.Block{
			Info:    mb.Title,
			Content: template.HTML(content), // #nosec G203 -- intended to be read from stdin

			Region:  mb.Region,
			BlockID: id,
		})

		if err != nil {
			_, _ = context.EPrintln(err.Error())
			return errBlocksGeneric
		}
	}

	return nil
}
