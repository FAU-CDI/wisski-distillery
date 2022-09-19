package instances

import (
	"errors"

	"github.com/FAU-CDI/wisski-distillery/internal/component"
	"github.com/FAU-CDI/wisski-distillery/internal/component/sql"
	"github.com/FAU-CDI/wisski-distillery/internal/component/triplestore"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/goprogram/exit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Instances manages multiple WissKI Instances.
type Instances struct {
	component.ComponentBase

	TS  *triplestore.Triplestore
	SQL *sql.SQL
}

func (Instances) Name() string {
	return "instances"
}

// ErrWissKINotFound is returned when a WissKI is not found
var ErrWissKINotFound = errors.New("WissKI not found")

var errSQL = exit.Error{
	Message:  "Unknown SQL Error %s",
	ExitCode: exit.ExitGeneric,
}

// WissKI returns the WissKI with the provided slug, if it exists.
// It the WissKI does not exist, returns ErrWissKINotFound.
func (instances *Instances) WissKI(slug string) (i WissKI, err error) {
	sql := instances.SQL
	if err := sql.WaitQueryTable(); err != nil {
		return i, err
	}

	table, err := sql.QueryTable(false, models.InstanceTable)
	if err != nil {
		return i, err
	}

	// find the instance by slug
	query := table.Where(&models.Instance{Slug: slug}).Find(&i.Instance)
	switch {
	case query.Error != nil:
		return i, errSQL.WithMessageF(query.Error)
	case query.RowsAffected == 0:
		return i, ErrWissKINotFound
	default:
		i.instances = instances
		return i, nil
	}
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
func (instances *Instances) All() ([]WissKI, error) {
	return instances.find(true, func(table *gorm.DB) *gorm.DB {
		return table
	})
}

// WissKIs returns the WissKI instances with the provides slugs.
// If a slug does not exist, it is omitted from the result.
func (instances *Instances) WissKIs(slugs ...string) ([]WissKI, error) {
	return instances.find(true, func(table *gorm.DB) *gorm.DB {
		return table.Where("slug IN ?", slugs)
	})
}

// Load is like All, except that when no slugs are provided, it calls All.
func (instances *Instances) Load(slugs ...string) ([]WissKI, error) {
	if len(slugs) == 0 {
		return instances.All()
	}
	return instances.WissKIs(slugs...)
}

// find finds instances based on the provided query
func (instances *Instances) find(order bool, query func(table *gorm.DB) *gorm.DB) (results []WissKI, err error) {
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
	results = make([]WissKI, len(bks))
	for i, bk := range bks {
		results[i].Instance = bk
		results[i].instances = instances
	}

	return results, nil
}
