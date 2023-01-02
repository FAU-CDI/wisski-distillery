package component

import (
	"context"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// UserDeleteHook represents a hook that is called just before a user is deleted
type UserDeleteHook interface {
	Component

	// OnUserDelete is called right before a user is deleted
	OnUserDelete(ctx context.Context, user *models.User) error
}
