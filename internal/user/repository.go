package user

import (
	"context"

	"github.com/Lidne/praktika_MAI/internal/models"
)

// UserRepository Sell
type UserRepository interface {
	Create(ctx context.Context, product *models.User) error
	Update(ctx context.Context, product *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	FindAll(ctx context.Context) ([]models.User, error)
	Delete(ctx context.Context, id int) error
}
