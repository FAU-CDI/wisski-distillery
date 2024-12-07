//spellchecker:words news
package news

//spellchecker:words context embed html template http strings time github wisski distillery internal component server assets handling templating wdlog yuin goldmark meta gmmeta parser golang slices
import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/assets"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/handling"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/server/templating"
	"github.com/FAU-CDI/wisski-distillery/internal/wdlog"
	"github.com/yuin/goldmark"
	gmmeta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/exp/slices"
)

type News struct {
	component.Base
	dependencies struct {
		Templating *templating.Templating
		Handling   *handling.Handling
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

	slices.SortFunc(items, func(a, b Item) int {
		return a.Date.Compare(b.Date)
	})

	return items, nil
}

//go:embed "news.html"
var newsHTML []byte
var newsTemplate = templating.Parse[newsContext](
	"news.html", newsHTML, nil,

	templating.Title("News"),
	templating.Assets(assets.AssetsDefault),
)

type newsContext struct {
	templating.RuntimeFlags
	Items []Item
}

var (
	menuNews = component.MenuItem{Title: "News", Path: "/news/"}
)

// HandleRoute returns the handler for the requested path
func (news *News) HandleRoute(ctx context.Context, path string) (http.Handler, error) {
	tpl := newsTemplate.Prepare(
		news.dependencies.Templating,
		templating.Crumbs(
			menuNews,
		),
	)

	items, itemsErr := Items()
	if itemsErr != nil {
		wdlog.Of(ctx).Error("Unable to load news items", "error", itemsErr)
	}

	return tpl.HTMLHandler(news.dependencies.Handling, func(r *http.Request) (nc newsContext, err error) {
		nc.Items, err = items, itemsErr
		return
	}), nil
}
