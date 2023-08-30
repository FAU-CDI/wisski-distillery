package extras

import (
	"context"
	"html/template"

	"github.com/FAU-CDI/wisski-distillery/internal/phpx"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"

	_ "embed"
)

type Blocks struct {
	ingredient.Base
	Dependencies struct {
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

// Create creates a new block with the given title and html content
func (blocks *Blocks) Create(ctx context.Context, server *phpx.Server, block Block) (err error) {
	err = blocks.Dependencies.PHP.ExecScript(ctx, server, nil, blocksPHP, "create_basic_block", block.Info, block.Content, block.Region, block.BlockID)
	return err
}

func (blocks *Blocks) GetFooterRegion(ctx context.Context, server *phpx.Server) (region string, err error) {
	err = blocks.Dependencies.PHP.ExecScript(ctx, server, &region, blocksPHP, "get_footer_region")
	return
}
