//spellchecker:words instances
package instances

//spellchecker:words context errors path filepath github wisski distillery internal component instances malt models gorm clause
import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"

	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/instances/malt"
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Instances manages multiple WissKI Instances.
type Instances struct {
	component.Base
	dependencies struct {
		Malt *malt.Malt
		SQL  *sql.SQL

		InstanceTable *sql.InstanceTable
	}
}

func (instances *Instances) Path() string {
	return filepath.Join(component.GetStill(instances).Config.Paths.Root, "instances")
}

// ErrWissKINotFound is returned when a WissKI is not found.
var ErrWissKINotFound = errors.New("WissKI not found")

// use uses the non-nil wisski instance with this instances.
func (instances *Instances) use(wisski *wisski.WissKI) {
	wisski.Malt = instances.dependencies.Malt
}

// WissKI returns the WissKI with the provided slug, if it exists.
// It the WissKI does not exist, returns an error wrapping [ErrWissKINotFound].
func (instances *Instances) WissKI(ctx context.Context, slug string) (wissKI *wisski.WissKI, err error) {
	if slug == "" {
		return nil, ErrWissKINotFound
	}

	if err := instances.dependencies.SQL.Wait(ctx); err != nil {
		return nil, fmt.Errorf("failed to wait for database: %w", err)
	}

	table, err := sql.OpenInterface[models.Instance](ctx, instances.dependencies.SQL, instances.dependencies.InstanceTable)
	if err != nil {
		return nil, fmt.Errorf("failed to open interface: %w", err)
	}

	wissKI = new(wisski.WissKI)

	wissKI.Instance, err = table.Where("slug = ?", slug).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %w", ErrWissKINotFound, err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find instance: %w", err)
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
	sql := instances.dependencies.SQL
	if err := sql.Wait(ctx); err != nil {
		return false, fmt.Errorf("failed to wait for database: %w", err)
	}

	table, err := sql.OpenTable(ctx, instances.dependencies.InstanceTable)
	if err != nil {
		return false, fmt.Errorf("failed to query table: %w", err)
	}

	query := table.Select("count(*) > 0").Where("slug = ?", slug).Find(&ok)
	if query.Error != nil {
		return false, fmt.Errorf("failed to count instances: %w", query.Error)
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

// find finds instances based on the provided query.
func (instances *Instances) find(ctx context.Context, order bool, query func(table *gorm.DB) *gorm.DB) (results []*wisski.WissKI, err error) {
	sql := instances.dependencies.SQL
	if err := sql.Wait(ctx); err != nil {
		return nil, fmt.Errorf("failed to wait for query table: %w", err)
	}

	// open the bookkeeping table
	table, err := sql.OpenTable(ctx, instances.dependencies.InstanceTable)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
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
		return nil, fmt.Errorf("failed to find instances: %w", find.Error)
	}

	// make proper instances
	results = make([]*wisski.WissKI, len(bks))
	for i, bk := range bks {
		results[i] = new(wisski.WissKI)
		results[i].Instance = bk
		instances.use(results[i])
	}

	return results, nil
}
