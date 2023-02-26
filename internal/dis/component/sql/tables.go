package sql

import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/pkglib/reflectx"
)

// This file defines additional tables used by multiple components

type InstanceTable struct {
	component.Base
}

var (
	_ component.Table = (*InstanceTable)(nil)
)

func (*InstanceTable) TableInfo() component.TableInfo {
	return component.TableInfo{
		Model: reflectx.TypeOf[models.Instance](),
		Name:  models.InstanceTable,
	}
}

type LockTable struct {
	component.Base
}

var (
	_ component.Table = (*LockTable)(nil)
)

func (*LockTable) TableInfo() component.TableInfo {
	return component.TableInfo{
		Model: reflectx.TypeOf[models.Lock](),
		Name:  models.LockTable,
	}
}
