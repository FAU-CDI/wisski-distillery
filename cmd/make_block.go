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
	"github.com/tkw1536/goprogram/exit"
)

// MakeBlock is the 'make_block' command.
var MakeBlock wisski_distillery.Command = makeBlock{}

type makeBlock struct {
	Title  string `short:"t" long:"title" description:"title of block to create"`
	Region string `short:"r" long:"region" description:"optional region to assign block to"`
	Footer bool   `short:"f" long:"footer" description:"create block in the footer region"`

	Positionals struct {
		Slug string `positional-arg-name:"SLUG" required:"1-1" description:"slug of instance to create legal block for"`
	} `positional-args:"true"`
}

var errFooterAndRegion = exit.Error{
	Message:  "`--footer` and `--region` provided",
	ExitCode: exit.ExitCommandArguments,
}

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

var errBlocksGeneric = exit.Error{
	Message:  "unable to create block",
	ExitCode: exit.ExitGeneric,
}

var errBlocksFooterFailed = exit.Error{
	Message:  "unable to determine footer block",
	ExitCode: exit.ExitGeneric,
}

var errBlocksNoFooter = exit.Error{
	Message:  "no footer known for region",
	ExitCode: exit.ExitGeneric,
}

var errBlocksNoContent = exit.Error{
	Message:  "unable to read content from standard input",
	ExitCode: exit.ExitCommandArguments,
}

func (mb makeBlock) Run(context wisski_distillery.Context) error {
	// get the wisski
	instance, err := context.Environment.Instances().WissKI(context.Context, mb.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errPathbuilderWissKI, err)
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
			context.EPrintln(err.Error())
			return errBlocksGeneric
		}
	}

	return nil
}
