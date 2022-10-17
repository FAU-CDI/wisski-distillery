package instances

import (
	"errors"
	"path/filepath"

	"github.com/FAU-CDI/wisski-distillery/internal/component"

	"github.com/FAU-CDI/wisski-distillery/internal/component/exporter/logger"
	"github.com/FAU-CDI/wisski-distillery/internal/component/meta"
	"github.com/FAU-CDI/wisski-distillery/internal/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/component/triplestore"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski"
	"github.com/tkw1536/goprogram/exit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Instances manages multiple WissKI Instances.
type Instances struct {
	component.ComponentBase

	TS          *triplestore.Triplestore
	SQL         *sql.SQL
	Meta        *meta.Meta
	ExporterLog *logger.Logger
}

func (instances *Instances) Path() string {
	return filepath.Join(instances.Still.Config.DeployRoot, "instances")
}

// ErrWissKINotFound is returned when a WissKI is not found
var ErrWissKINotFound = errors.New("WissKI not found")

var errSQL = exit.Error{
	Message:  "Unknown SQL Error %s",
	ExitCode: exit.ExitGeneric,
}

// use uses the non-nil wisski instance with this instances
func (instances *Instances) use(wisski *wisski.WissKI) {
	wisski.Core = instances.Still
	wisski.SQL = instances.SQL
	wisski.TS = instances.TS
	wisski.Meta = instances.Meta
	wisski.ExporterLog = instances.ExporterLog
}

// WissKI returns the WissKI with the provided slug, if it exists.
// It the WissKI does not exist, returns ErrWissKINotFound.
func (instances *Instances) WissKI(slug string) (wissKI *wisski.WissKI, err error) {
	sql := instances.SQL
	if err := sql.WaitQueryTable(); err != nil {
		return nil, err
	}

	table, err := sql.QueryTable(false, models.InstanceTable)
	if err != nil {
		return nil, err
	}

	// create a struct
	wissKI = new(wisski.WissKI)

	// find the instance by slug
	query := table.Where(&models.Instance{Slug: slug}).Find(&wissKI.Instance)
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
func (instances *Instances) Instance(instance models.Instance) *wisski.WissKI {
	wissKI, err := instances.WissKI(instance.Slug)
	if err != nil {
		return nil
	}
	return wissKI
}

// Has checks if a WissKI with the provided slug exists inside the database.
// It does not perform any checks on the WissKI itself.
func (instances *Instances) Has(slug string) (ok bool, err error) {
	sql := instances.SQL
	if err := sql.WaitQueryTable(); err != nil {
		return false, err
	}

	table, err := sql.QueryTable(false, models.InstanceTable)
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
func (instances *Instances) All() ([]*wisski.WissKI, error) {
	return instances.find(true, func(table *gorm.DB) *gorm.DB {
		return table
	})
}

// WissKIs returns the WissKI instances with the provides slugs.
// If a slug does not exist, it is omitted from the result.
func (instances *Instances) WissKIs(slugs ...string) ([]*wisski.WissKI, error) {
	return instances.find(true, func(table *gorm.DB) *gorm.DB {
		return table.Where("slug IN ?", slugs)
	})
}

// Load is like All, except that when no slugs are provided, it calls All.
func (instances *Instances) Load(slugs ...string) ([]*wisski.WissKI, error) {
	if len(slugs) == 0 {
		return instances.All()
	}
	return instances.WissKIs(slugs...)
}

// find finds instances based on the provided query
func (instances *Instances) find(order bool, query func(table *gorm.DB) *gorm.DB) (results []*wisski.WissKI, err error) {
	sql := instances.SQL
	if err := sql.WaitQueryTable(); err != nil {
		return nil, err
	}

	// open the bookkeeping table
	table, err := sql.QueryTable(false, models.InstanceTable)
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
		results[i].Instance = bk
		instances.use(results[i])
	}

	return results, nil
}
