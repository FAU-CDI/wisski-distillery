package instances

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/malt"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/goprogram/exit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Instances manages multiple WissKI Instances.
type Instances struct {
	component.Base
	Dependencies struct {
		Malt *malt.Malt
		SQL  *sql.SQL
	}
}

func (instances *Instances) Path() string {
	return filepath.Join(instances.Still.Config.DeployRoot, "instances")
}

// ErrWissKINotFound is returned when a WissKI is not found
var ErrWissKINotFound = errors.New("WissKI not found")

var errSQL = exit.Error{
	Message:  "unknown SQL error %s",
	ExitCode: exit.ExitGeneric,
}

// use uses the non-nil wisski instance with this instances
func (instances *Instances) use(wisski *wisski.WissKI) {
	wisski.Liquid.Malt = instances.Dependencies.Malt
}

// WissKI returns the WissKI with the provided slug, if it exists.
// It the WissKI does not exist, returns ErrWissKINotFound.
func (instances *Instances) WissKI(ctx context.Context, slug string) (wissKI *wisski.WissKI, err error) {
	sql := instances.Dependencies.SQL
	if err := sql.WaitQueryTable(ctx); err != nil {
		return nil, err
	}

	table, err := sql.QueryTable(ctx, false, models.InstanceTable)
	if err != nil {
		return nil, err
	}

	// create a struct
	wissKI = new(wisski.WissKI)

	// find the instance by slug
	query := table.Where(&models.Instance{Slug: slug}).Find(&wissKI.Liquid.Instance)
	switch {
	case query.Error != nil:
		return nil, errSQL.WithMessageF(query.Error)
	case query.RowsAffected == 0:
		return nil, ErrWissKINotFound
	}

	// use the wissKI instance
	instances.use(wissKI)
	return wissKI, nil
}

// Instance is a convenience function to return an instance based on a model slug.
// When the instance does not exist, returns nil.
func (instances *Instances) Instance(ctx context.Context, instance models.Instance) *wisski.WissKI {
	wissKI, err := instances.WissKI(ctx, instance.Slug)
	if err != nil {
		return nil
	}
	return wissKI
}

// Has checks if a WissKI with the provided slug exists inside the database.
// It does not perform any checks on the WissKI itself.
func (instances *Instances) Has(ctx context.Context, slug string) (ok bool, err error) {
	sql := instances.Dependencies.SQL
	if err := sql.WaitQueryTable(ctx); err != nil {
		return false, err
	}

	table, err := sql.QueryTable(ctx, false, models.InstanceTable)
	if err != nil {
		return false, err
	}

	query := table.Select("count(*) > 0").Where("slug = ?", slug).Find(&ok)
	if query.Error != nil {
		return false, errSQL.WithMessageF(query.Error)
	}
	return
}

// All returns all instances of the WissKI Distillery in consistent order.
//
// There is no guarantee that this order remains identical between different api releases; however subsequent invocations are guaranteed to return the same order.
func (instances *Instances) All(ctx context.Context) ([]*wisski.WissKI, error) {
	return instances.find(ctx, true, func(table *gorm.DB) *gorm.DB {
		return table
	})
}

// WissKIs returns the WissKI instances with the provides slugs.
// If a slug does not exist, it is omitted from the result.
func (instances *Instances) WissKIs(ctx context.Context, slugs ...string) ([]*wisski.WissKI, error) {
	return instances.find(ctx, true, func(table *gorm.DB) *gorm.DB {
		return table.Where("slug IN ?", slugs)
	})
}

// Load is like All, except that when no slugs are provided, it calls All.
func (instances *Instances) Load(ctx context.Context, slugs ...string) ([]*wisski.WissKI, error) {
	if len(slugs) == 0 {
		return instances.All(ctx)
	}
	return instances.WissKIs(ctx, slugs...)
}

// find finds instances based on the provided query
func (instances *Instances) find(ctx context.Context, order bool, query func(table *gorm.DB) *gorm.DB) (results []*wisski.WissKI, err error) {
	sql := instances.Dependencies.SQL
	if err := sql.WaitQueryTable(ctx); err != nil {
		return nil, err
	}

	// open the bookkeeping table
	table, err := sql.QueryTable(ctx, false, models.InstanceTable)
	if err != nil {
		return nil, err
	}

	// prepare a query
	find := table
	if order {
		find = find.Order(clause.OrderByColumn{Column: clause.Column{Name: "slug"}, Desc: false})
	}
	if query != nil {
		find = query(find)
	}

	// fetch bookkeeping instances
	var bks []models.Instance
	find = find.Find(&bks)
	if find.Error != nil {
		return nil, errSQL.WithMessageF(find.Error)
	}

	// make proper instances
	results = make([]*wisski.WissKI, len(bks))
	for i, bk := range bks {
		results[i] = new(wisski.WissKI)
		results[i].Liquid.Instance = bk
		instances.use(results[i])
	}

	return results, nil
}
