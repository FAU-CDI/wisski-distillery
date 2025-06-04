package sql

//spellchecker:words reflect github wisski distillery internal component models
import (
	"github.com/FAU-CDI/wisski-distillery/internal/dis/component"
	"github.com/FAU-CDI/wisski-distillery/internal/models"
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
		Model: models.Instance{},
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
		Model: models.Lock{},
	}
}
