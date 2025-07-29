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
	"github.com/spf13/cobra"
	"go.tkw01536.de/pkglib/exit"
)

func NewMakeBlockCommand() *cobra.Command {
	impl := new(makeBlock)

	cmd := &cobra.Command{
		Use:     "make_block",
		Short:   "Creates a block with html content provided on stdin",
		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.StringVar(&impl.Title, "title", "", "title of block to create")
	flags.StringVar(&impl.Region, "region", "", "optional region to assign block to")
	flags.BoolVar(&impl.Footer, "footer", false, "create block in the footer region")

	return cmd
}

type makeBlock struct {
	Title       string
	Region      string
	Footer      bool
	Positionals struct {
		Slug string
	}
}

func (mb *makeBlock) ParseArgs(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		mb.Positionals.Slug = args[0]
	}

	if mb.Region != "" && mb.Footer {
		return errFooterAndRegion
	}
	return nil
}

func (*makeBlock) Description() wisski_distillery.Description {
	return wisski_distillery.Description{
		Requirements: cli.Requirements{
			NeedsDistillery: true,
		},
		Command:     "make_block",
		Description: "Creates a block with html content provided on stdin",
	}
}

var errFooterAndRegion = exit.NewErrorWithCode("`--footer` and `--region` provided", exit.ExitCommandArguments)
var (
	errBlocksGeneric      = exit.NewErrorWithCode("unable to create block", exit.ExitGeneric)
	errBlocksFooterFailed = exit.NewErrorWithCode("unable to determine footer block", exit.ExitGeneric)
	errBlocksNoFooter     = exit.NewErrorWithCode("no footer known for region", exit.ExitGeneric)
	errBlocksNoContent    = exit.NewErrorWithCode("unable to read content from standard input", exit.ExitCommandArguments)
)

func (mb *makeBlock) Exec(cmd *cobra.Command, args []string) error {
	dis, err := cli.GetDistillery(cmd, cli.Requirements{
		NeedsDistillery: true,
	})

	if err != nil {
		return fmt.Errorf("%w: %w", errPathbuildersWissKI, err)
	}

	// get the wisski
	instance, err := dis.Instances().WissKI(cmd.Context(), mb.Positionals.Slug)
	if err != nil {
		return fmt.Errorf("%w: %w", errPathbuildersWissKI, err)
	}

	// get the footer (if any)
	if mb.Footer {
		wdlog.Of(cmd.Context()).Info("checking for footer")
		region, err := instance.Blocks().GetFooterRegion(cmd.Context(), nil)
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
	content, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return fmt.Errorf("%w: %w", errBlocksNoContent, err)
	}

	{
		err := instance.Blocks().Create(cmd.Context(), nil, extras.Block{
			Info:    mb.Title,
			Content: template.HTML(content), // #nosec G203 -- intended to be read from stdin

			Region:  mb.Region,
			BlockID: id,
		})

		if err != nil {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err.Error())
			return errBlocksGeneric
		}
	}

	return nil
}
