//spellchecker:words extras
package extras

//spellchecker:words context html template github wisski distillery internal phpx ingredient embed
import (
	"context"
	"fmt"
	"html/template"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"

	_ "embed"
)

type Blocks struct {
	ingredient.Base
	dependencies struct {
		PHP *php.PHP
	}
}

//go:embed blocks.php
var blocksPHP string

type Block struct {
	Info    string
	Content template.HTML

	Region  string
	BlockID string
}

// Create creates a new block with the given title and html content.
func (blocks *Blocks) Create(ctx context.Context, server *phpx.Server, block Block) error {
	err := blocks.dependencies.PHP.ExecScript(ctx, server, nil, blocksPHP, "create_basic_block", block.Info, block.Content, block.Region, block.BlockID)
	if err != nil {
		return fmt.Errorf("failed to create block: %w", err)
	}
	return nil
}

func (blocks *Blocks) GetFooterRegion(ctx context.Context, server *phpx.Server) (region string, err error) {
	err = blocks.dependencies.PHP.ExecScript(ctx, server, &region, blocksPHP, "get_footer_region")
	if err != nil {
		return "", fmt.Errorf("failed to get footer region: %w", err)
	}
	return region, nil
}
