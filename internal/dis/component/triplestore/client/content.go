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

const NQuadsContentType = "application/n-quads"

var (
	errExportWrongStatusCode  = errors.New("ExportContent: Wrong status code")
	errReplaceWrongStatusCode = errors.New("ReplaceContent: Wrong status code")
)

// ExportContent exports the content of the provided repository as an n-quads file and writes them into dst.
// count contains the total number of bytes written, and any error.
func (client *Client) ExportContent(ctx context.Context, dst io.Writer, repo string) (c int64, e error) {
	res, err := client.rest(ctx, http.MethodGet, "/repositories/"+url.PathEscape(repo)+"/statements?infer=false", headers{Accept: NQuadsContentType})
	if err != nil {
		return 0, fmt.Errorf("failed to send statements endpoint request: %w", err)
	}
	defer errorsx.Close(res.Body, &e, "response body")

	if err := newStatusError(res, true, http.StatusOK); err != nil {
		return 0, fmt.Errorf("statements endpoint responded: %w", err)
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
	res, err := client.doRestWithReader(ctx, http.MethodPut, "/repositories/"+url.PathEscape(repo)+"/statements", headers{ContentType: NQuadsContentType}, reader)
	if err != nil {
		return fmt.Errorf("failed to send statements endpoint request: %w", err)
	}
	defer func() {
		// we don't care about any errors of closing the body
		_ = res.Body.Close()
	}()

	if err := newStatusError(res, true, http.StatusNoContent); err != nil {
		return fmt.Errorf("statements endpoint responded: %w", err)
	}
	return nil
}
