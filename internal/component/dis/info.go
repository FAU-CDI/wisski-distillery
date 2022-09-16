package dis

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/bookkeeping"
	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/FAU-CDI/wisski-distillery/internal/config"
	"github.com/FAU-CDI/wisski-distillery/pkg/httpx"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/sync/errgroup"

	_ "embed"
)

func (dis *Dis) info(io stream.IOStream) (http.Handler, error) {
	mux := http.NewServeMux()

	// handle everything under /dis/!
	mux.HandleFunc("/dis/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/dis/" {
			http.Redirect(w, r, "/dis/index", http.StatusTemporaryRedirect)
			return
		}
		http.NotFound(w, r)
	})

	// static stuff
	static, err := dis.disStatic()
	if err != nil {
		return nil, err
	}
	mux.Handle("/dis/static/", static)

	// render everything
	mux.Handle("/dis/index", httpx.HTMLHandler[disIndex]{
		Handler:  dis.disIndex,
		Template: indexTemplate,
	})

	mux.Handle("/dis/instance/", httpx.HTMLHandler[disInstance]{
		Handler:  dis.disInstance,
		Template: instanceTemplate,
	})

	// api -- for future usage
	mux.Handle("/dis/api/v1/instance/get/", httpx.JSON(dis.getinstance))
	mux.Handle("/dis/api/v1/instance/all", httpx.JSON(dis.allinstances))

	// ensure that everyone is logged in!
	return httpx.BasicAuth(mux, "WissKI Distillery Admin", func(user, pass string) bool {
		return user == dis.Config.DisAdminUser && pass == dis.Config.DisAdminPassword
	}), nil
}

// disIndex is the context of the "/dis/index" page
type disIndex struct {
	Time time.Time

	Config *config.Config

	Instances    []instances.Info
	TotalCount   int
	RunningCount int
	StoppedCount int
}

func (dis *Dis) disIndex(r *http.Request) (idx disIndex, err error) {
	// load instances
	idx.Instances, err = dis.allinstances(r)
	if err != nil {
		return
	}

	// count how many are running and how many are stopped
	for _, i := range idx.Instances {
		if i.Running {
			idx.RunningCount++
		} else {
			idx.StoppedCount++
		}
	}
	idx.TotalCount = len(idx.Instances)

	// get the static properties
	idx.Config = dis.Config

	// current time
	idx.Time = time.Now()

	return
}

// disInstance is the context of the "/dis/instance/*" page
type disInstance struct {
	Time time.Time

	Instance bookkeeping.Instance
	Info     instances.Info
}

func (dis *Dis) disInstance(r *http.Request) (is disInstance, err error) {
	// find the slug as the last component of path!
	slug := strings.TrimSuffix(r.URL.Path, "/")
	slug = slug[strings.LastIndex(slug, "/")+1:]

	// find the instance itself!
	instance, err := dis.Instances.WissKI(slug)
	if err == instances.ErrWissKINotFound {
		return is, httpx.ErrNotFound
	}
	if err != nil {
		return is, err
	}
	is.Instance = instance.Instance

	// get some more info about the wisski
	is.Info, err = instance.Info(false)
	if err != nil {
		return is, err
	}

	// current time
	is.Time = time.Now()

	return
}

//go:embed html/static
var htmlStaticFS embed.FS

func (*Dis) disStatic() (http.Handler, error) {
	fs, err := fs.Sub(htmlStaticFS, "html/static")
	if err != nil {
		return nil, err
	}

	return http.StripPrefix("/dis/static/", http.FileServer(http.FS(fs))), nil
}

//go:embed "html/index.html"
var indexTemplateStr string
var indexTemplate = template.Must(template.New("index.html").Parse(indexTemplateStr))

//go:embed "html/instance.html"
var instanceTemplateString string
var instanceTemplate = template.Must(template.New("instance.html").Parse(instanceTemplateString))

func (dis *Dis) getinstance(r *http.Request) (info instances.Info, err error) {
	// find the slug as the last component of path!
	slug := strings.TrimSuffix(r.URL.Path, "/")
	slug = slug[strings.LastIndex(slug, "/")+1:]

	// load the wisski instance!
	wisski, err := dis.Instances.WissKI(strings.TrimSuffix(slug, "/"))
	if err == instances.ErrWissKINotFound {
		return info, httpx.ErrNotFound
	}
	if err != nil {
		return info, err
	}

	// get info about it!
	return wisski.Info(false)
}

func (dis *Dis) allinstances(*http.Request) (infos []instances.Info, err error) {
	var errgroup errgroup.Group

	// list all the instances
	all, err := dis.Instances.All()
	if err != nil {
		return nil, err
	}

	// get all of their info!
	infos = make([]instances.Info, len(all))
	for i, instance := range all {
		{
			i := i
			instance := instance

			errgroup.Go(func() (err error) {
				infos[i], err = instance.Info(true)
				return err
			})
		}
	}

	// wait for the results, and return
	err = errgroup.Wait()
	return
}
