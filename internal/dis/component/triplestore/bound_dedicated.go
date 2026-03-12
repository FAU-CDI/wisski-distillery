package triplestore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/triplestore/client"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/pkg/dockerx"
	"go.tkw01536.de/pkglib/stream"
	"go.tkw01536.de/pkglib/timex"
)

// Template for creating repositories.
//
// NOTE(twiesing): The template is not aware of SparQL syntax, thus this template is very unsafe.
// And should only be used with KNOWN GOOD input.
var createRepoTemplate = template.Must(template.New("bound_dedicated_createrepo.tpl").Parse(`
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#>.
@prefix config: <tag:rdf4j.org,2023:config/>.

[] a config:Repository ;
   config:rep.id "{{ .RepositoryID }}" ;
   rdfs:label "{{ .Label }}" ;
   config:rep.impl [
      config:rep.type "openrdf:SailRepository" ;
      config:sail.impl [
        config:sail.type "openrdf:NativeStore" ;
         config:sail.iterationCacheSyncThreshold "10000";
         config:sail.defaultQueryEvaluationMode "STANDARD";
         config:native.tripleIndexes "spoc,posc"
      ]
   ].
`))

// boundDedicated implements a wrapper around the dedicated triplestore client.
type boundDedicated struct {
	openStack   func() (*dockerx.Stack, error)
	serviceName string

	instance models.Instance
}

func (bound *boundDedicated) ReadURL() string {
	return "http://dedicatedtriplestore:8080/rdf4j-server/repositories/" + url.PathEscape(bound.instance.GraphDBRepository)
}

func (bound *boundDedicated) WriteURL() string {
	return "http://dedicatedtriplestore:8080/rdf4j-server/repositories/" + url.PathEscape(bound.instance.GraphDBRepository) + "/statements"
}

func (bound *boundDedicated) Credentials() (username string, password string) {
	return "", ""
}

// RestoreDB snapshots the provided repository into dst.
func (bound *boundDedicated) RestoreDB(ctx context.Context, progress io.Writer, reader io.Reader) (e error) {
	return bound.do(ctx, stream.Null, true, func(stack *dockerx.Stack) error {
		return bound.curl(
			ctx, stack, progress,
			"PUT", "/repositories/"+url.PathEscape(bound.instance.GraphDBRepository)+"/statements",
			map[string]string{"Content-Type": client.NQuadsContentType},
			reader,
		)
	})
}

// Purge purges the given repository.
func (bound *boundDedicated) Purge(ctx context.Context, progress io.Writer, allowCreate bool) error {
	return bound.do(ctx, progress, allowCreate, func(stack *dockerx.Stack) error {
		return bound.curl(ctx, stack, stream.Null, "DELETE", "/repositories/"+url.PathEscape(bound.instance.GraphDBRepository), nil, nil)
	})
}

// SnapshotDB snapshots the provided repository into dst.
func (bound *boundDedicated) SnapshotDB(ctx context.Context, progress io.Writer, dst io.Writer) error {
	return bound.do(ctx, stream.Null, true, func(stack *dockerx.Stack) error {
		return bound.curl(
			ctx, stack, dst,
			"GET", "/repositories/"+url.PathEscape(bound.instance.GraphDBRepository)+"/statements?infer=false",
			map[string]string{"Accept": client.NQuadsContentType},
			nil,
		)
	})
}

// Provision provisions the repository for this instance, possibly deleting any existing repositories.
func (bound *boundDedicated) Provision(ctx context.Context, progress io.Writer, domain string) (e error) {
	var createRepo bytes.Buffer
	if err := createRepoTemplate.Execute(&createRepo, client.CreateOpts{
		RepositoryID: bound.instance.GraphDBRepository,
		Label:        domain,
		BaseURL:      "http://" + domain + "/",
	}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return bound.do(ctx, progress, true, func(stack *dockerx.Stack) error {
		return bound.curl(ctx, stack, stream.Null, "PUT", "/repositories/"+url.PathEscape(bound.instance.GraphDBRepository), map[string]string{"Content-Type": "text/turtle"}, &createRepo)
	})
}

func (bound *boundDedicated) do(ctx context.Context, progress io.Writer, allowCreate bool, fn func(stack *dockerx.Stack) error) (e error) {
	if err := dockerx.Do(ctx, stream.Null, allowCreate, bound.openStack, func(stack *dockerx.Stack) error {
		if err := timex.TickUntilFunc(func(time.Time) bool {
			return bound.curl(ctx, stack, progress, "GET", "/protocol", nil, nil) == nil
		}, ctx, time.Second); err != nil {
			return fmt.Errorf("failed to wait for triplestore to be ready: %w", err)
		}

		return fn(stack)
	}, bound.serviceName); err != nil {
		return fmt.Errorf("dockerx.Do returned: %w", err)
	}
	return nil
}

var errNonZeroExitCode = errors.New("non-zero exit code")

// curl executes a curl request against the dedicated triplestore.
func (bound *boundDedicated) curl(ctx context.Context, stack *dockerx.Stack, stdout io.Writer, method string, path string, headers map[string]string, stdin io.Reader) (e error) {
	command := makeCurlCommand(method, "http://localhost:8080/rdf4j-server"+path, headers, true, stdin != nil)

	var errBuf bytes.Buffer

	if code := stack.Exec(ctx, stream.NewIOStream(stdout, &errBuf, stdin), dockerx.ExecOptions{
		Service: bound.serviceName,
		Cmd:     command[0],
		Args:    command[1:],
	})(); code != 0 {
		return fmt.Errorf("%w: %s returned non-zero exit code: %d: %s", errNonZeroExitCode, strings.Join(command, " "), code, errBuf.String())
	}
	return nil
}

// makeCurlCommand generates a curl command for the given method, url and headers.
func makeCurlCommand(method string, url string, headers map[string]string, fail bool, stdin bool) []string {
	command := []string{
		"curl",
		"--no-progress-meter",
		"--verbose",
		"--request", method,
	}
	for key, value := range headers {
		command = append(command, "--header", fmt.Sprintf("%s: %s", key, value))
	}
	if fail {
		command = append(command, "--fail")
	}
	if stdin {
		command = append(command, "--data-binary", "@-")
	}
	command = append(command, url)
	return command
}
