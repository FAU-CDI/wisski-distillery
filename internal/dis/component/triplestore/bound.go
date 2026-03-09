package triplestore

import (
	"context"
	"io"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"go.tkw01536.de/pkglib/errorsx"
)

// TODO: Move these into a bound struct

// RestoreDB snapshots the provided repository into dst.
func (ts Triplestore) RestoreDB(ctx context.Context, repo string, reader io.Reader) (e error) {
	return ts.client().ReplaceContent(ctx, repo, reader)
}

// Purge purges the given repository and user.
func (ts *Triplestore) Purge(ctx context.Context, instance models.Instance, domain string) error {
	client := ts.client()
	return errorsx.Combine(
		client.DeleteRepository(ctx, instance.GraphDBRepository),
		client.DeleteUser(ctx, instance.GraphDBUsername),
	)
}

// SnapshotDB snapshots the provided repository into dst.
func (ts Triplestore) SnapshotDB(ctx context.Context, dst io.Writer, repo string) (c int64, e error) {
	return ts.client().ExportContent(ctx, dst, repo)
}
