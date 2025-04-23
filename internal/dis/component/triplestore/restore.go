//spellchecker:words triplestore
package triplestore

//spellchecker:words context http github errors
import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var errTSRestoreWrongStatusCode = errors.New("Triplestore.Restore: Wrong status code")

// RestoreDB snapshots the provided repository into dst.
func (ts Triplestore) RestoreDB(ctx context.Context, repo string, reader io.Reader) error {
	// submit the form
	res, err := ts.DoRestWithReader(ctx, 0, http.MethodPut, "/repositories/"+repo+"/statements", &RequestHeaders{ContentType: nquadsContentType}, reader)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		message, _ := io.ReadAll(res.Body)
		return fmt.Errorf("%w: %s", errTSRestoreWrongStatusCode, message)
	}
	return nil
}
