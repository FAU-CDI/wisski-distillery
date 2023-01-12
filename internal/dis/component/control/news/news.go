package news

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/control/static/custom"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/rs/zerolog"
	"github.com/yuin/goldmark"
	gmmeta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/exp/slices"
)

type News struct {
	component.Base
	Dependencies struct {
		Custom *custom.Custom
	}
}

var (
	_ component.Routeable = (*News)(nil)
)

func (*News) Routes() component.Routes {
	return component.Routes{
		Prefix: "/news/",
		Exact:  true,
		CSRF:   false,

		MenuTitle:    "News",
		MenuPriority: component.MenuNews,
	}
}

type Item struct {
	ID      string
	Date    time.Time
	Title   string
	Content template.HTML
}

func (item *Item) parse(path string, builder *strings.Builder) error {
	builder.Reset()

	// open file
	content, err := fs.ReadFile(newsFS, path)
	if err != nil {
		return err
	}

	// parse and read metadata
	reader := goldmark.New(goldmark.WithExtensions(
		gmmeta.Meta,
	))

	context := parser.NewContext()
	if err := reader.Convert(content, builder, parser.WithContext(context)); err != nil {
		return err
	}
	meta := gmmeta.Get(context)

	// read title
	item.Title, _ = meta["title"].(string)

	// read date
	date, _ := meta["date"].(string)
	item.Date, err = time.Parse("2006-01-02", date)
	if err != nil {
		return err
	}

	// write content
	item.Content = template.HTML(builder.String())

	return nil
}

//go:embed "NEWS/*.md"
var newsFS embed.FS

// Items returns a list of all news items
func Items() ([]Item, error) {
	var builder strings.Builder

	files, err := fs.Glob(newsFS, "NEWS/*.md")
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(files))
	for i, file := range files {
		items[i].ID = file[len("NEWS/") : len(file)-len(".md")]
		if err := items[i].parse(file, &builder); err != nil {
			return nil, err
		}
	}

	slices.SortFunc(items, func(a, b Item) bool {
		return !a.Date.Before(b.Date)
	})

	return items, nil
}

//go:embed "news.html"
var newsHTMLStr string
var newsTemplate = static.AssetsDefault.MustParseShared("news.html", newsHTMLStr)

type newsContext struct {
	custom.BaseContext
	Items []Item
}

// HandleRoute returns the handler for the requested path
func (news *News) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	gaps := custom.BaseContextGaps{
		Crumbs: []component.MenuItem{
			{Title: "News", Path: "/news/"},
		},
	}

	items, itemsErr := Items()
	if itemsErr != nil {
		zerolog.Ctx(ctx).Err(itemsErr).Msg("Unable to load news items")
	}

	return httpx.HTMLHandler[newsContext]{
		Handler: func(r *http.Request) (nc newsContext, err error) {
			news.Dependencies.Custom.Update(&nc, r, gaps)
			nc.Items, err = items, itemsErr

			return
		},
		Template: news.Dependencies.Custom.Template(newsTemplate),
	}, nil
}
