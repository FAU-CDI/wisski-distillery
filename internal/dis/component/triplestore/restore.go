//spellchecker:words triplestore
package triplestore

//spellchecker:words context http github errors
import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tkw1536/pkglib/errorsx"
)

var errTSRestoreWrongStatusCode = errors.New("Triplestore.Restore: Wrong status code")

// RestoreDB snapshots the provided repository into dst.
func (ts Triplestore) RestoreDB(ctx context.Context, repo string, reader io.Reader) (e error) {
	// submit the form
	res, err := ts.DoRestWithReader(ctx, 0, http.MethodPut, "/repositories/"+url.PathEscape(repo)+"/statements", &RequestHeaders{ContentType: nquadsContentType}, reader)
	if err != nil {
		return err
	}
	defer errorsx.Close(res.Body, &e, "request body")

	if res.StatusCode != http.StatusNoContent {
		message, _ := io.ReadAll(res.Body)
		return fmt.Errorf("%w: %s", errTSRestoreWrongStatusCode, message)
	}
	return nil
}
