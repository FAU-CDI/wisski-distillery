package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.tkw01536.de/pkglib/errorsx"
)

const nquadsContentType = "text/x-nquads"

var (
	errExportWrongStatusCode  = errors.New("ExportContent: Wrong status code")
	errReplaceWrongStatusCode = errors.New("ReplaceContent: Wrong status code")
)

// ExportContent exports the content of the provided repository as an n-quads file and writes them into dst.
// count contains the total number of bytes written, and any error.
func (client *Client) ExportContent(ctx context.Context, dst io.Writer, repo string) (c int64, e error) {
	res, err := client.rest(ctx, http.MethodGet, "/repositories/"+url.PathEscape(repo)+"/statements?infer=false", &requestHeaders{Accept: nquadsContentType})
	if err != nil {
		return 0, fmt.Errorf("failed to send rest request: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	if res.StatusCode != http.StatusOK {
		return 0, errExportWrongStatusCode
	}
	count, err := io.Copy(dst, res.Body)
	if err != nil {
		return count, fmt.Errorf("failed to copy result: %w", err)
	}
	return count, nil
}

// ReplaceContent repleaces the content of the provided repository with the content of the given reader.
// The reader must contain valid n-quads data.
func (client *Client) ReplaceContent(ctx context.Context, repo string, reader io.Reader) (e error) {
	res, err := client.doRestWithReader(ctx, http.MethodPut, "/repositories/"+url.PathEscape(repo)+"/statements", &requestHeaders{ContentType: nquadsContentType}, reader)
	if err != nil {
		return err
	}
	defer func() {
		// we don't care about any errors of closing the body
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusNoContent {
		message, _ := io.ReadAll(res.Body)
		return fmt.Errorf("%w: %s", errReplaceWrongStatusCode, message)
	}
	return nil
}
